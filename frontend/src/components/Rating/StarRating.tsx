import React from 'react';
import { FaStar, FaStarHalfAlt } from 'react-icons/fa';
import './StarRating.css';

interface StarRatingProps {
  rating: number;
  maxRating?: number;
  interactive?: boolean;
  onRatingChange?: (rating: number) => void;
  size?: 'small' | 'medium' | 'large';
}

export const StarRating: React.FC<StarRatingProps> = ({
  rating,
  maxRating = 10,
  interactive = false,
  onRatingChange,
  size = 'medium',
}) => {
  const starRating = (rating / maxRating) * 5;
  const fullStars = Math.floor(starRating);
  const hasHalfStar = starRating % 1 >= 0.5;
  const emptyStars = 5 - fullStars - (hasHalfStar ? 1 : 0);

  const handleStarClick = (index: number) => {
    if (interactive && onRatingChange) {
      const newRating = ((index + 1) / 5) * maxRating;
      onRatingChange(newRating);
    }
  };

  const sizeClass = `star-rating-${size}`;

  return (
    <div className={`star-rating ${sizeClass} ${interactive ? 'interactive' : ''}`}>
      {[...Array(fullStars)].map((_, i) => (
        <FaStar
          key={i}
          className="star star-full"
          onClick={() => handleStarClick(i)}
        />
      ))}
      {hasHalfStar && (
        <FaStarHalfAlt
          className="star star-half"
          onClick={() => handleStarClick(fullStars)}
        />
      )}
      {[...Array(emptyStars)].map((_, i) => (
        <FaStar
          key={i + fullStars + (hasHalfStar ? 1 : 0)}
          className="star star-empty"
          onClick={() => handleStarClick(i + fullStars + (hasHalfStar ? 1 : 0))}
        />
      ))}
      {!interactive && <span className="rating-value">{rating.toFixed(1)}</span>}
    </div>
  );
};

