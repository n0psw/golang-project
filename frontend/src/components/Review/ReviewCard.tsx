import React from 'react';
import { type Review } from '../../api/client';
import { StarRating } from '../Rating/StarRating';
import { formatDate } from '../../utils/helpers';
import './ReviewCard.css';

interface ReviewCardProps {
  review: Review;
  onEdit?: () => void;
  onDelete?: () => void;
  canEdit?: boolean;
}

export const ReviewCard: React.FC<ReviewCardProps> = ({
  review,
  onEdit,
  onDelete,
  canEdit = false,
}) => {
  return (
    <div className="review-card">
      <div className="review-card-header">
        <div className="review-card-author">
          <span className="review-username">{review.user?.username || 'Неизвестный пользователь'}</span>
          <span className="review-date">{formatDate(review.created_at)}</span>
        </div>
        <div className="review-card-rating">
          <StarRating rating={review.rating} size="small" />
        </div>
      </div>
      <h3 className="review-card-title">{review.title}</h3>
      <p className="review-card-content">{review.content}</p>
      {canEdit && (
        <div className="review-card-actions">
          {onEdit && (
            <button className="review-action-btn edit" onClick={onEdit}>
              Редактировать
            </button>
          )}
          {onDelete && (
            <button className="review-action-btn delete" onClick={onDelete}>
              Удалить
            </button>
          )}
        </div>
      )}
    </div>
  );
};

