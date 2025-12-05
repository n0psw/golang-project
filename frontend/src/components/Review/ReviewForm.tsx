import React, { useState, useEffect } from 'react';
import { type Review, type CreateReviewRequest } from '../../api/client';
import { Input } from '../common/Input';
import { StarRating } from '../Rating/StarRating';
import { Button } from '../common/Button';
import './ReviewForm.css';

interface ReviewFormProps {
  review?: Review;
  onSubmit: (data: CreateReviewRequest) => Promise<void>;
  onCancel: () => void;
}

export const ReviewForm: React.FC<ReviewFormProps> = ({
  review,
  onSubmit,
  onCancel,
}) => {
  const [formData, setFormData] = useState<CreateReviewRequest>({
    rating: 5,
    title: '',
    content: '',
  });
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (review) {
      setFormData({
        rating: review.rating,
        title: review.title,
        content: review.content,
      });
    } else {
      setFormData({
        rating: 5,
        title: '',
        content: '',
      });
    }
  }, [review]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
    if (errors[name]) {
      setErrors((prev) => ({ ...prev, [name]: '' }));
    }
  };

  const handleRatingChange = (rating: number) => {
    setFormData((prev) => ({ ...prev, rating }));
  };

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {};
    if (!formData.title.trim()) {
      newErrors.title = 'Заголовок обязателен';
    }
    if (!formData.content.trim()) {
      newErrors.content = 'Текст отзыва обязателен';
    }
    if (formData.rating < 1 || formData.rating > 10) {
      newErrors.rating = 'Рейтинг должен быть от 1 до 10';
    }
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;
    setLoading(true);
    try {
      await onSubmit(formData);
    } catch (error) {
      console.error('Error submitting form:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="review-form">
      <div className="review-form-rating">
        <label className="form-label">Рейтинг</label>
        <StarRating
          rating={formData.rating}
          interactive={true}
          onRatingChange={handleRatingChange}
          size="large"
        />
        {errors.rating && (
          <span className="input-error-text">{errors.rating}</span>
        )}
      </div>
      <Input
        label="Заголовок"
        name="title"
        value={formData.title}
        onChange={handleChange}
        error={errors.title}
        required
      />
      <div className="form-group">
        <label className="form-label">
          Текст отзыва
          <textarea
            name="content"
            value={formData.content}
            onChange={handleChange}
            className="form-textarea"
            rows={6}
            required
          />
        </label>
        {errors.content && (
          <span className="input-error-text">{errors.content}</span>
        )}
      </div>
      <div className="form-actions">
        <Button type="submit" disabled={loading}>
          {review ? 'Сохранить' : 'Опубликовать'}
        </Button>
        <Button type="button" variant="outline" onClick={onCancel}>
          Отмена
        </Button>
      </div>
    </form>
  );
};

