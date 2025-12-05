import React, { useState, useEffect } from 'react';
import { useSearchParams, Link } from 'react-router-dom';
import { moviesAPI, genresAPI, reviewsAPI, type Movie, type Genre, type Review, type CreateMovieRequest } from '../api/client';
import { useAuth } from '../context/AuthContext';
import { MovieForm } from '../components/Movie/MovieForm';
import { Modal } from '../components/common/Modal';
import { Button } from '../components/common/Button';
import { Input } from '../components/common/Input';
import { Loading } from '../components/common/Loading';
import { StarRating } from '../components/Rating/StarRating';
import { formatDate } from '../utils/helpers';
import './AdminPage.css';

export const AdminPage: React.FC = () => {
  const { isAdmin } = useAuth();
  const [searchParams] = useSearchParams();
  const [activeTab, setActiveTab] = useState<'movies' | 'genres' | 'reviews'>('movies');
  const [movies, setMovies] = useState<Movie[]>([]);
  const [genres, setGenres] = useState<Genre[]>([]);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [loading, setLoading] = useState(true);
  const [reviewsLoading, setReviewsLoading] = useState(false);
  const [reviewsCurrentPage, setReviewsCurrentPage] = useState(1);
  const [reviewsTotalPages, setReviewsTotalPages] = useState(1);
  const [showMovieForm, setShowMovieForm] = useState(false);
  const [showGenreForm, setShowGenreForm] = useState(false);
  const [editingMovie, setEditingMovie] = useState<Movie | null>(null);
  const [newGenreName, setNewGenreName] = useState('');

  useEffect(() => {
    if (!isAdmin()) return;
    const editMovieId = searchParams.get('editMovie');
    if (!editMovieId) return;
    const movie = movies.find((m) => m.id === editMovieId);
    if (movie) {
      setEditingMovie(movie);
      setShowMovieForm(true);
    }
  }, [searchParams, isAdmin, movies]);

  useEffect(() => {
    if (!isAdmin()) return;
    loadMovies();
    loadGenres();
  }, [isAdmin]);

  useEffect(() => {
    if (!isAdmin()) return;
    if (activeTab === 'reviews') {
      loadReviews();
    }
  }, [isAdmin, activeTab, reviewsCurrentPage]);

  const loadMovies = async () => {
    try {
      const response = await moviesAPI.getAll({ limit: 1000 });
      const payload = response.data;
      const moviesData = Array.isArray(payload?.data) ? (payload.data as Movie[]) : [];
      const moviesArray = moviesData.map((movie: Movie) => ({
        ...movie,
        genres: Array.isArray(movie.genres) ? movie.genres : [],
      }));
      setMovies(moviesArray);
    } catch (error) {
      console.error('Error loading movies:', error);
      setMovies([]);
    } finally {
      setLoading(false);
    }
  };

  const loadGenres = async () => {
    try {
      const response = await genresAPI.getAll();
      const genresData = Array.isArray(response.data) ? response.data : [];
      setGenres(genresData);
    } catch (error) {
      console.error('Error loading genres:', error);
      setGenres([]);
    }
  };

  const loadReviews = async () => {
    setReviewsLoading(true);
    try {
      const response = await reviewsAPI.getAll({
        page: reviewsCurrentPage,
        limit: 20,
      });
      const payload = response.data;
      const reviewsData = Array.isArray(payload?.data)
        ? (payload.data as Review[])
        : Array.isArray(payload?.reviews)
          ? (payload.reviews as Review[])
          : Array.isArray(payload)
            ? (payload as Review[])
            : [];
      const reviewsArray = reviewsData;
      setReviews(reviewsArray);
      setReviewsTotalPages(payload?.total_pages ?? 1);
    } catch (error) {
      console.error('Error loading reviews:', error);
      setReviews([]);
    } finally {
      setReviewsLoading(false);
    }
  };

  const handleCreateMovie = async (data: CreateMovieRequest) => {
    try {
      const requestData = {
        ...data,
        duration_minutes: data.duration,
        genre_ids: data.genre_ids,
      };
      await moviesAPI.create(requestData);
      setShowMovieForm(false);
      setEditingMovie(null);
      loadMovies();
    } catch (error) {
      console.error('Error creating movie:', error);
      throw error;
    }
  };

  const handleUpdateMovie = async (data: CreateMovieRequest) => {
    if (!editingMovie) return;
    try {
      const requestData = {
        ...data,
        duration_minutes: data.duration,
        genre_ids: data.genre_ids,
      };
      await moviesAPI.update(editingMovie.id, requestData);
      setShowMovieForm(false);
      setEditingMovie(null);
      loadMovies();
    } catch (error) {
      console.error('Error updating movie:', error);
      throw error;
    }
  };

  const handleDeleteMovie = async (movie: Movie) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот фильм?')) return;
    try {
      await moviesAPI.delete(movie.id);
      loadMovies();
    } catch (error) {
      console.error('Error deleting movie:', error);
    }
  };

  const handleCreateGenre = async () => {
    if (!newGenreName.trim()) return;
    try {
      await genresAPI.create({ name: newGenreName.trim() });
      setNewGenreName('');
      setShowGenreForm(false);
      loadGenres();
    } catch (error) {
      console.error('Error creating genre:', error);
    }
  };

  const handleUpdateGenre = async (id: string, name: string) => {
    try {
      await genresAPI.update(id, { name });
      loadGenres();
    } catch (error) {
      console.error('Error updating genre:', error);
    }
  };

  const handleDeleteGenre = async (genre: Genre) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот жанр?')) return;
    try {
      await genresAPI.delete(genre.id);
      loadGenres();
    } catch (error) {
      console.error('Error deleting genre:', error);
    }
  };

  const handleDeleteReview = async (review: Review) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот отзыв?')) return;
    try {
      await reviewsAPI.delete(review.id);
      loadReviews();
    } catch (error) {
      console.error('Error deleting review:', error);
    }
  };

  if (!isAdmin()) {
    return (
      <div className="page">
        <div className="container">
          <div className="admin-error">
            <p>Доступ запрещен. Требуются права администратора.</p>
          </div>
        </div>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="page">
        <div className="container">
          <Loading />
        </div>
      </div>
    );
  }

  return (
    <div className="page">
      <div className="container">
        <h1 className="page-title">Админ панель</h1>
        <div className="admin-tabs">
          <button
            className={`admin-tab ${activeTab === 'movies' ? 'active' : ''}`}
            onClick={() => setActiveTab('movies')}
          >
            Фильмы
          </button>
          <button
            className={`admin-tab ${activeTab === 'genres' ? 'active' : ''}`}
            onClick={() => setActiveTab('genres')}
          >
            Жанры
          </button>
          <button
            className={`admin-tab ${activeTab === 'reviews' ? 'active' : ''}`}
            onClick={() => setActiveTab('reviews')}
          >
            Отзывы
          </button>
        </div>
        {activeTab === 'movies' && (
          <div className="admin-section">
            <div className="admin-section-header">
              <h2>Управление фильмами</h2>
              <Button onClick={() => {
                setEditingMovie(null);
                setShowMovieForm(true);
              }}>
                Добавить фильм
              </Button>
            </div>
            <div className="admin-movies-list">
              {movies && Array.isArray(movies) && movies.map((movie) => (
                <div key={movie.id} className="admin-movie-item">
                  <div className="admin-movie-info">
                    <h3>{movie.title}</h3>
                    <p>{movie.release_year} • {movie.director}</p>
                  </div>
                  <div className="admin-movie-actions">
                    <Button
                      variant="outline"
                      onClick={() => {
                        setEditingMovie(movie);
                        setShowMovieForm(true);
                      }}
                    >
                      Редактировать
                    </Button>
                    <Button
                      variant="danger"
                      onClick={() => handleDeleteMovie(movie)}
                    >
                      Удалить
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
        {activeTab === 'genres' && (
          <div className="admin-section">
            <div className="admin-section-header">
              <h2>Управление жанрами</h2>
              <Button onClick={() => setShowGenreForm(true)}>
                Добавить жанр
              </Button>
            </div>
            <div className="admin-genres-list">
              {genres && Array.isArray(genres) && genres.map((genre) => (
                <div key={genre.id} className="admin-genre-item">
                  <Input
                    value={genre.name}
                    onChange={(e) => {
                      const newName = e.target.value;
                      if (newName !== genre.name) {
                        handleUpdateGenre(genre.id, newName);
                      }
                    }}
                    className="admin-genre-input"
                  />
                  <Button
                    variant="danger"
                    onClick={() => handleDeleteGenre(genre)}
                  >
                    Удалить
                  </Button>
                </div>
              ))}
            </div>
          </div>
        )}
        {activeTab === 'reviews' && (
          <div className="admin-section">
            <div className="admin-section-header">
              <h2>Управление отзывами</h2>
            </div>
            {reviewsLoading ? (
              <Loading />
            ) : reviews.length === 0 ? (
              <div className="admin-reviews-empty">
                <p>Отзывы отсутствуют</p>
              </div>
            ) : (
              <div className="admin-reviews-list">
                {reviews && Array.isArray(reviews) && reviews.map((review) => (
                  <div key={review.id} className="admin-review-item">
                    <div className="admin-review-header">
                      <div className="admin-review-meta">
                        <Link to={`/movie/${review.movie_id}`} className="admin-review-movie-link">
                          {review.movie?.title || `Фильм #${review.movie_id}`}
                        </Link>
                        <span className="admin-review-author">от {review.user?.username || 'Неизвестный пользователь'}</span>
                        <span className="admin-review-date">{formatDate(review.created_at)}</span>
                      </div>
                      <div className="admin-review-rating">
                        <StarRating rating={review.rating} size="small" />
                      </div>
                    </div>
                    <h3 className="admin-review-title">{review.title}</h3>
                    <p className="admin-review-content">{review.content}</p>
                    <div className="admin-review-actions">
                      <Button
                        variant="danger"
                        onClick={() => handleDeleteReview(review)}
                      >
                        Удалить
                      </Button>
                    </div>
                  </div>
                ))}
                {reviewsTotalPages > 1 && (
                  <div className="admin-reviews-pagination">
                    <button
                      className="pagination-btn"
                      onClick={() => {
                        setReviewsCurrentPage(reviewsCurrentPage - 1);
                      }}
                      disabled={reviewsCurrentPage === 1}
                    >
                      Назад
                    </button>
                    <span className="pagination-info">
                      Страница {reviewsCurrentPage} из {reviewsTotalPages}
                    </span>
                    <button
                      className="pagination-btn"
                      onClick={() => {
                        setReviewsCurrentPage(reviewsCurrentPage + 1);
                      }}
                      disabled={reviewsCurrentPage === reviewsTotalPages}
                    >
                      Вперед
                    </button>
                  </div>
                )}
              </div>
            )}
          </div>
        )}
        <Modal
          isOpen={showMovieForm}
          onClose={() => {
            setShowMovieForm(false);
            setEditingMovie(null);
          }}
          title={editingMovie ? 'Редактировать фильм' : 'Добавить фильм'}
        >
          <MovieForm
            movie={editingMovie || undefined}
            genres={genres}
            onSubmit={editingMovie ? handleUpdateMovie : handleCreateMovie}
            onCancel={() => {
              setShowMovieForm(false);
              setEditingMovie(null);
            }}
          />
        </Modal>
        <Modal
          isOpen={showGenreForm}
          onClose={() => {
            setShowGenreForm(false);
            setNewGenreName('');
          }}
          title="Добавить жанр"
        >
          <div className="genre-form">
            <Input
              label="Название жанра"
              value={newGenreName}
              onChange={(e) => setNewGenreName(e.target.value)}
              placeholder="Введите название жанра"
            />
            <div className="form-actions">
              <Button onClick={handleCreateGenre} disabled={!newGenreName.trim()}>
                Создать
              </Button>
              <Button variant="outline" onClick={() => {
                setShowGenreForm(false);
                setNewGenreName('');
              }}>
                Отмена
              </Button>
            </div>
          </div>
        </Modal>
      </div>
    </div>
  );
};
