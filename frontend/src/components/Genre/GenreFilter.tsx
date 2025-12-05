import React from 'react';
import { type Genre } from '../../api/client';
import './GenreFilter.css';

interface GenreFilterProps {
  genres: Genre[];
  selectedGenreId: string | null;
  onGenreChange: (genreId: string | null) => void;
}

export const GenreFilter: React.FC<GenreFilterProps> = ({
  genres,
  selectedGenreId,
  onGenreChange,
}) => {
  if (!genres || !Array.isArray(genres)) {
    return null;
  }

  return (
    <div className="genre-filter">
      <label className="genre-filter-label">Жанр:</label>
      <select
        className="genre-filter-select"
        value={selectedGenreId || ''}
        onChange={(e) => onGenreChange(e.target.value || null)}
      >
        <option value="">Все жанры</option>
        {genres.map((genre) => (
          <option key={genre.id} value={genre.id}>
            {genre.name}
          </option>
        ))}
      </select>
    </div>
  );
};

