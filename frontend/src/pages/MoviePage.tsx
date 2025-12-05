import React, { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { moviesAPI, reviewsAPI, type Movie, type Review, type CreateReviewRequest } from '../api/client';
import { useAuth } from '../context/AuthContext';
import { MovieDetail } from '../components/Movie/MovieDetail';
import { ReviewList } from '../components/Review/ReviewList';
import { ReviewForm } from '../components/Review/ReviewForm';
import { Modal } from '../components/common/Modal';
import { Button } from '../components/common/Button';
import { Loading } from '../components/common/Loading';
import './MoviePage.css';

export const MoviePage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user, isAdmin } = useAuth();
  const [movie, setMovie] = useState<Movie | null>(null);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [loading, setLoading] = useState(true);
  const [reviewsLoading, setReviewsLoading] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [showReviewForm, setShowReviewForm] = useState(false);
  const [editingReview, setEditingReview] = useState<Review | null>(null);
  const [userReview, setUserReview] = useState<Review | null>(null);

  useEffect(() => {
    const loadMovie = async () => {
      if (!id) return;
      setLoading(true);
      try {
        const response = await moviesAPI.getById(id);
        const movieData = response.data;
        if (movieData && movieData.genres && !Array.isArray(movieData.genres)) {
          movieData.genres = [];
        }
        setMovie(movieData);
      } catch (error) {
        console.error('Error loading movie:', error);
        navigate('/');
      } finally {
        setLoading(false);
      }
    };
    loadMovie();
  }, [id, navigate]);

  const fetchReviews = useCallback(async () => {
    if (!id) return;
    setReviewsLoading(true);
    try {
      const response = await reviewsAPI.getByMovie(id, {
        page: currentPage,
        limit: 10,
      });
      const payload = response.data;
      const reviewsData = Array.isArray(payload?.data) ? (payload.data as Review[]) : [];
      setReviews(reviewsData);
      setTotalPages(payload?.total_pages ?? 1);

      if (user) {
        const userReviewEntry = reviewsData.find((r) => r.user_id === user.id) || null;
        setUserReview(userReviewEntry);
      } else {
        setUserReview(null);
      }
    } catch (error) {
      console.error('Error loading reviews:', error);
      setReviews([]);
      setUserReview(null);
    } finally {
      setReviewsLoading(false);
    }
  }, [id, currentPage, user?.id]);

  useEffect(() => {
    fetchReviews();
  }, [fetchReviews]);

  const handleCreateReview = async (data: CreateReviewRequest) => {
    if (!movie) {
      return;
    }
    try {
      await reviewsAPI.create(movie.id, data);
      setShowReviewForm(false);
      setEditingReview(null);
      await fetchReviews();
    } catch (error) {
      console.error('Error creating review:', error);
      throw error;
    }
  };

  const handleUpdateReview = async (data: CreateReviewRequest) => {
    if (!editingReview) return;
    try {
      await reviewsAPI.update(editingReview.id, data);
      setShowReviewForm(false);
      setEditingReview(null);
      await fetchReviews();
    } catch (error) {
      console.error('Error updating review:', error);
      throw error;
    }
  };

  const handleDeleteReview = async (review: Review) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот отзыв?')) return;
    try {
      await reviewsAPI.delete(review.id);
      await fetchReviews();
    } catch (error) {
      console.error('Error deleting review:', error);
    }
  };

  const handleDeleteMovie = async () => {
    if (!movie) return;
    if (!window.confirm('Вы уверены, что хотите удалить этот фильм?')) return;
    try {
      await moviesAPI.delete(movie.id);
      navigate('/');
    } catch (error) {
      console.error('Error deleting movie:', error);
    }
  };

  const canEditReview = (review: Review): boolean => {
    return !!(user && (review.user_id === user.id || isAdmin()));
  };

  if (loading || !movie) {
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
        <MovieDetail
          movie={movie}
          onEdit={() => navigate(`/admin?editMovie=${movie.id}`)}
          onDelete={handleDeleteMovie}
          canEdit={isAdmin()}
        />
        <div className="movie-page-reviews">
          <div className="reviews-header">
            <h2 className="reviews-title">Отзывы</h2>
            {user && !userReview && (
              <Button onClick={() => setShowReviewForm(true)}>
                Написать отзыв
              </Button>
            )}
          </div>
          <ReviewList
            reviews={reviews}
            loading={reviewsLoading}
            currentPage={currentPage}
            totalPages={totalPages}
            onPageChange={setCurrentPage}
            onEdit={(review) => {
              setEditingReview(review);
              setShowReviewForm(true);
            }}
            onDelete={handleDeleteReview}
            canEdit={canEditReview}
          />
        </div>
        <Modal
          isOpen={showReviewForm}
          onClose={() => {
            setShowReviewForm(false);
            setEditingReview(null);
          }}
          title={editingReview ? 'Редактировать отзыв' : 'Написать отзыв'}
        >
          <ReviewForm
            review={editingReview || undefined}
            onSubmit={editingReview ? handleUpdateReview : handleCreateReview}
            onCancel={() => {
              setShowReviewForm(false);
              setEditingReview(null);
            }}
          />
        </Modal>
      </div>
    </div>
  );
};
