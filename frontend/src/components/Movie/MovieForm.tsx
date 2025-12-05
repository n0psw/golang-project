import React, { useState, useEffect } from 'react';
import { type Movie, type Genre, type CreateMovieRequest } from '../../api/client';
import { Input } from '../common/Input';
import { Button } from '../common/Button';
import './MovieForm.css';

interface MovieFormProps {
  movie?: Movie;
  genres: Genre[];
  onSubmit: (data: CreateMovieRequest) => Promise<void>;
  onCancel: () => void;
}

export const MovieForm: React.FC<MovieFormProps> = ({
  movie,
  genres,
  onSubmit,
  onCancel,
}) => {
  const [formData, setFormData] = useState<CreateMovieRequest>({
    title: '',
    release_year: new Date().getFullYear(),
    director: '',
    duration: 90,
    description: '',
    genre_ids: [],
  });
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);
  const [submitError, setSubmitError] = useState<string>('');

  useEffect(() => {
    if (movie) {
      setFormData({
        title: movie.title,
        release_year: movie.release_year,
        director: movie.director,
        duration: movie.duration_minutes || movie.duration,
        description: movie.description,
        genre_ids: movie.genres && Array.isArray(movie.genres) ? movie.genres.map((g) => g.id) : [],
      });
    }
  }, [movie]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: name === 'release_year' || name === 'duration' ? Number(value) : value,
    }));
    if (errors[name]) {
      setErrors((prev) => ({ ...prev, [name]: '' }));
    }
  };

  const handleGenreToggle = (genreId: string) => {
    setFormData((prev) => ({
      ...prev,
      genre_ids: prev.genre_ids.includes(genreId)
        ? prev.genre_ids.filter((id) => id !== genreId)
        : [...prev.genre_ids, genreId],
    }));
  };

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {};
    if (!formData.title.trim()) {
      newErrors.title = 'Название обязательно';
    }
    if (formData.release_year < 1900 || formData.release_year > new Date().getFullYear()) {
      newErrors.release_year = 'Некорректный год';
    }
    if (!formData.director.trim()) {
      newErrors.director = 'Режиссер обязателен';
    }
    if (formData.duration < 1) {
      newErrors.duration = 'Длительность должна быть больше 0';
    }
    if (formData.genre_ids.length === 0) {
      newErrors.genre_ids = 'Выберите хотя бы один жанр';
    }
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;
    setLoading(true);
    setSubmitError('');
    try {
      await onSubmit(formData);
    } catch (error: any) {
      console.error('Error submitting form:', error);
      const errorMessage = error?.response?.data?.error || error?.message || 'Ошибка при создании фильма';
      setSubmitError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="movie-form">
      {submitError && (
        <div className="form-error" style={{ color: 'red', marginBottom: '1rem', padding: '0.75rem', backgroundColor: '#fee', borderRadius: '4px' }}>
          {submitError}
        </div>
      )}
      <Input
        label="Название"
        name="title"
        value={formData.title}
        onChange={handleChange}
        error={errors.title}
        required
      />
      <div className="form-row">
        <Input
          type="number"
          label="Год выпуска"
          name="release_year"
          value={formData.release_year.toString()}
          onChange={handleChange}
          error={errors.release_year}
          required
        />
        <Input
          type="number"
          label="Длительность (мин)"
          name="duration"
          value={formData.duration.toString()}
          onChange={handleChange}
          error={errors.duration}
          required
        />
      </div>
      <Input
        label="Режиссер"
        name="director"
        value={formData.director}
        onChange={handleChange}
        error={errors.director}
        required
      />
      <div className="form-group">
        <label className="form-label">
          Описание
          <textarea
            name="description"
            value={formData.description}
            onChange={handleChange}
            className="form-textarea"
            rows={5}
          />
        </label>
      </div>
      <div className="form-group">
        <label className="form-label">Жанры</label>
        <div className="genre-checkboxes">
          {genres && Array.isArray(genres) && genres.length > 0 ? (
            genres.map((genre) => (
              <label key={genre.id} className="genre-checkbox">
                <input
                  type="checkbox"
                  checked={formData.genre_ids.includes(genre.id)}
                  onChange={() => handleGenreToggle(genre.id)}
                />
                <span>{genre.name}</span>
              </label>
            ))
          ) : (
            <p className="genre-loading">Загрузка жанров...</p>
          )}
        </div>
        {errors.genre_ids && (
          <span className="input-error-text">{errors.genre_ids}</span>
        )}
      </div>
      <div className="form-actions">
        <Button type="submit" disabled={loading}>
          {movie ? 'Сохранить' : 'Создать'}
        </Button>
        <Button type="button" variant="outline" onClick={onCancel}>
          Отмена
        </Button>
      </div>
    </form>
  );
};

