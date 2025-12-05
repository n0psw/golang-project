import React from 'react';
import { Link } from 'react-router-dom';
import { type Movie } from '../../api/client';
import { StarRating } from '../Rating/StarRating';
import { GenreBadge } from '../Genre/GenreBadge';
import './MovieCard.css';

interface MovieCardProps {
  movie: Movie;
}

export const MovieCard: React.FC<MovieCardProps> = ({ movie }) => {
  return (
    <Link to={`/movie/${movie.id}`} className="movie-card">
      <div className="movie-card-poster">
        <div className="movie-card-placeholder">
          {movie.title.charAt(0).toUpperCase()}
        </div>
      </div>
      <div className="movie-card-content">
        <h3 className="movie-card-title">{movie.title}</h3>
        <p className="movie-card-year">{movie.release_year}</p>
        <div className="movie-card-rating">
          <StarRating rating={movie.average_rating} size="small" />
        </div>
        <div className="movie-card-genres">
          {movie.genres && Array.isArray(movie.genres) && movie.genres.slice(0, 3).map((genre) => (
            <GenreBadge key={genre.id} name={genre.name} />
          ))}
        </div>
      </div>
    </Link>
  );
};

