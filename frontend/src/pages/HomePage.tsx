import React, { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { moviesAPI, genresAPI, type Movie, type Genre } from '../api/client';
import { MovieFilters } from '../components/Filters/MovieFilters';
import { MovieList } from '../components/Movie/MovieList';
import './HomePage.css';

export const HomePage: React.FC = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const [movies, setMovies] = useState<Movie[]>([]);
  const [genres, setGenres] = useState<Genre[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);

  const selectedGenreId = searchParams.get('genre') || null;
  const year = searchParams.get('year') || '';
  const minRating = searchParams.get('minRating')
    ? Number(searchParams.get('minRating'))
    : 0;
  const search = searchParams.get('search') || '';

  useEffect(() => {
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
    loadGenres();
  }, []);

  useEffect(() => {
    const loadMovies = async () => {
      setLoading(true);
      try {
        const params: Record<string, string | number> = {
          page: currentPage,
          limit: 12,
        };
        if (selectedGenreId) params.genre_id = selectedGenreId;
        if (year) params.year = Number(year);
        if (minRating > 0) params.min_rating = minRating;
        if (search) params.search = search;

        const response = await moviesAPI.getAll(params);
        const payload = response.data;
        const moviesData = Array.isArray(payload?.data) ? (payload.data as Movie[]) : [];
        const moviesArray = moviesData.map((movie: Movie) => ({
          ...movie,
          genres: Array.isArray(movie.genres) ? movie.genres : [],
        }));
        setMovies(moviesArray);
        setTotalPages(payload?.total_pages ?? 1);
      } catch (error) {
        console.error('Error loading movies:', error);
        setMovies([]);
      } finally {
        setLoading(false);
      }
    };
    loadMovies();
  }, [currentPage, selectedGenreId, year, minRating, search]);

  const handleGenreChange = (genreId: string | null) => {
    setSearchParams((prev) => {
      if (genreId) {
        prev.set('genre', genreId);
      } else {
        prev.delete('genre');
      }
      prev.delete('page');
      return prev;
    });
    setCurrentPage(1);
  };

  const handleYearChange = (newYear: string) => {
    setSearchParams((prev) => {
      if (newYear) {
        prev.set('year', newYear);
      } else {
        prev.delete('year');
      }
      prev.delete('page');
      return prev;
    });
    setCurrentPage(1);
  };

  const handleMinRatingChange = (rating: number) => {
    setSearchParams((prev) => {
      if (rating > 0) {
        prev.set('minRating', rating.toString());
      } else {
        prev.delete('minRating');
      }
      prev.delete('page');
      return prev;
    });
    setCurrentPage(1);
  };

  const handleSearchChange = (newSearch: string) => {
    setSearchParams((prev) => {
      if (newSearch) {
        prev.set('search', newSearch);
      } else {
        prev.delete('search');
      }
      prev.delete('page');
      return prev;
    });
    setCurrentPage(1);
  };

  const handleReset = () => {
    setSearchParams({});
    setCurrentPage(1);
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    setSearchParams((prev) => {
      prev.set('page', page.toString());
      return prev;
    });
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  return (
    <div className="page">
      <div className="container">
        <h1 className="page-title">Фильмы</h1>
        <MovieFilters
          genres={genres}
          selectedGenreId={selectedGenreId}
          onGenreChange={handleGenreChange}
          year={year}
          onYearChange={handleYearChange}
          minRating={minRating}
          onMinRatingChange={handleMinRatingChange}
          search={search}
          onSearchChange={handleSearchChange}
          onReset={handleReset}
        />
        <MovieList
          movies={movies}
          loading={loading}
          currentPage={currentPage}
          totalPages={totalPages}
          onPageChange={handlePageChange}
        />
      </div>
    </div>
  );
};
