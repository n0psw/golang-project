import React from 'react';
import { type Movie } from '../../api/client';
import { MovieCard } from './MovieCard';
import { Loading } from '../common/Loading';
import './MovieList.css';

interface MovieListProps {
  movies: Movie[];
  loading: boolean;
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

export const MovieList: React.FC<MovieListProps> = ({
  movies,
  loading,
  currentPage,
  totalPages,
  onPageChange,
}) => {
  if (loading) {
    return <Loading />;
  }

  if (!movies || !Array.isArray(movies) || movies.length === 0) {
    return (
      <div className="movie-list-empty">
        <p>Фильмы не найдены</p>
      </div>
    );
  }

  return (
    <div className="movie-list">
      <div className="movie-list-grid">
        {movies.map((movie) => (
          <MovieCard key={movie.id} movie={movie} />
        ))}
      </div>
      {totalPages > 1 && (
        <div className="movie-list-pagination">
          <button
            className="pagination-btn"
            onClick={() => onPageChange(currentPage - 1)}
            disabled={currentPage === 1}
          >
            Назад
          </button>
          <span className="pagination-info">
            Страница {currentPage} из {totalPages}
          </span>
          <button
            className="pagination-btn"
            onClick={() => onPageChange(currentPage + 1)}
            disabled={currentPage === totalPages}
          >
            Вперед
          </button>
        </div>
      )}
    </div>
  );
};

