import React from 'react';
import { type Movie } from '../../api/client';
import { StarRating } from '../Rating/StarRating';
import { GenreBadge } from '../Genre/GenreBadge';
import './MovieDetail.css';

interface MovieDetailProps {
  movie: Movie;
  onEdit?: () => void;
  onDelete?: () => void;
  canEdit?: boolean;
}

export const MovieDetail: React.FC<MovieDetailProps> = ({
  movie,
  onEdit,
  onDelete,
  canEdit = false,
}) => {
  return (
    <div className="movie-detail">
      <div className="movie-detail-header">
        <div className="movie-detail-poster">
          <div className="movie-detail-placeholder">
            {movie.title.charAt(0).toUpperCase()}
          </div>
        </div>
        <div className="movie-detail-info">
          <h1 className="movie-detail-title">{movie.title}</h1>
          <div className="movie-detail-meta">
            <span>{movie.release_year}</span>
            <span>•</span>
            <span>{movie.director}</span>
            <span>•</span>
            <span>{movie.duration_minutes || movie.duration} мин</span>
          </div>
          <div className="movie-detail-rating">
            <StarRating rating={movie.average_rating} size="large" />
          </div>
          <div className="movie-detail-genres">
            {movie.genres && Array.isArray(movie.genres) && movie.genres.map((genre) => (
              <GenreBadge key={genre.id} name={genre.name} />
            ))}
          </div>
          {canEdit && (
            <div className="movie-detail-actions">
              {onEdit && (
                <button className="action-btn edit-btn" onClick={onEdit}>
                  Редактировать
                </button>
              )}
              {onDelete && (
                <button className="action-btn delete-btn" onClick={onDelete}>
                  Удалить
                </button>
              )}
            </div>
          )}
        </div>
      </div>
      <div className="movie-detail-description">
        <h2>Описание</h2>
        <p>{movie.description || 'Описание отсутствует'}</p>
      </div>
    </div>
  );
};

