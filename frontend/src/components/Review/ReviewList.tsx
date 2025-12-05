import React from 'react';
import { type Review } from '../../api/client';
import { ReviewCard } from './ReviewCard';
import { Loading } from '../common/Loading';
import './ReviewList.css';

interface ReviewListProps {
  reviews: Review[];
  loading: boolean;
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  onEdit?: (review: Review) => void;
  onDelete?: (review: Review) => void;
  canEdit?: (review: Review) => boolean;
}

export const ReviewList: React.FC<ReviewListProps> = ({
  reviews,
  loading,
  currentPage,
  totalPages,
  onPageChange,
  onEdit,
  onDelete,
  canEdit,
}) => {
  if (loading) {
    return <Loading />;
  }

  if (!reviews || !Array.isArray(reviews) || reviews.length === 0) {
    return (
      <div className="review-list-empty">
        <p>Отзывы отсутствуют</p>
      </div>
    );
  }

  return (
    <div className="review-list">
      {reviews.map((review) => (
        <ReviewCard
          key={review.id}
          review={review}
          onEdit={onEdit ? () => onEdit(review) : undefined}
          onDelete={onDelete ? () => onDelete(review) : undefined}
          canEdit={canEdit ? canEdit(review) : false}
        />
      ))}
      {totalPages > 1 && (
        <div className="review-list-pagination">
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

