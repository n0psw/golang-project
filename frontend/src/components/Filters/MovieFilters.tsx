import React from 'react';
import { type Genre } from '../../api/client';
import { GenreFilter } from '../Genre/GenreFilter';
import { Input } from '../common/Input';
import { Button } from '../common/Button';
import './MovieFilters.css';

interface MovieFiltersProps {
  genres: Genre[];
  selectedGenreId: string | null;
  onGenreChange: (genreId: string | null) => void;
  year: string;
  onYearChange: (year: string) => void;
  minRating: number;
  onMinRatingChange: (rating: number) => void;
  search: string;
  onSearchChange: (search: string) => void;
  onReset: () => void;
}

export const MovieFilters: React.FC<MovieFiltersProps> = ({
  genres,
  selectedGenreId,
  onGenreChange,
  year,
  onYearChange,
  minRating,
  onMinRatingChange,
  search,
  onSearchChange,
  onReset,
}) => {
  return (
    <div className="movie-filters">
      <div className="filters-row">
        <Input
          type="text"
          placeholder="Поиск по названию..."
          value={search}
          onChange={(e) => onSearchChange(e.target.value)}
          className="filter-search"
        />
        <GenreFilter
          genres={genres}
          selectedGenreId={selectedGenreId}
          onGenreChange={onGenreChange}
        />
        <div className="filter-group">
          <label className="filter-label">Год:</label>
          <Input
            type="number"
            placeholder="Год"
            value={year}
            onChange={(e) => onYearChange(e.target.value)}
            className="filter-year"
          />
        </div>
        <div className="filter-group">
          <label className="filter-label">Мин. рейтинг: {minRating}</label>
          <input
            type="range"
            min="0"
            max="10"
            step="0.5"
            value={minRating}
            onChange={(e) => onMinRatingChange(Number(e.target.value))}
            className="filter-slider"
          />
        </div>
        <Button onClick={onReset} variant="outline">
          Сбросить
        </Button>
      </div>
    </div>
  );
};

