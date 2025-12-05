import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { reviewsAPI, type Review } from '../api/client';
import { ReviewCard } from '../components/Review/ReviewCard';
import { Loading } from '../components/common/Loading';
import './MyReviewsPage.css';

export const MyReviewsPage: React.FC = () => {
  const [reviews, setReviews] = useState<Review[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);

  useEffect(() => {
    const loadReviews = async () => {
      setLoading(true);
      try {
        const response = await reviewsAPI.getMyReviews({
          page: currentPage,
          limit: 10,
        });
        const payload = response.data;
        const reviewsData = Array.isArray(payload?.data)
          ? (payload.data as Review[])
          : Array.isArray(payload?.reviews)
            ? (payload.reviews as Review[])
            : Array.isArray(payload)
              ? (payload as Review[])
              : [];
        const reviewsArray = reviewsData;
        setReviews(reviewsArray);
        setTotalPages(payload?.total_pages ?? 1);
      } catch (error) {
        console.error('Error loading reviews:', error);
        setReviews([]);
      } finally {
        setLoading(false);
      }
    };
    loadReviews();
  }, [currentPage]);

  const handleDeleteReview = async (review: Review) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот отзыв?')) return;
    try {
      await reviewsAPI.delete(review.id);
      setReviews((prev) => prev.filter((r) => r.id !== review.id));
    } catch (error) {
      console.error('Error deleting review:', error);
    }
  };

  if (loading && reviews.length === 0) {
    return (
      <div className="page">
        <div className="container">
          <Loading />
        </div>
      </div>
    );
  }

  return (
    <div className="page">
      <div className="container">
        <h1 className="page-title">Мои отзывы</h1>
        {reviews.length === 0 && !loading ? (
          <div className="my-reviews-empty">
            <p>У вас пока нет отзывов</p>
            <Link to="/" className="link-button">
              Перейти к фильмам
            </Link>
          </div>
        ) : (
          <div className="my-reviews-list">
            {reviews && Array.isArray(reviews) && reviews.map((review) => (
              <div key={review.id} className="my-review-item">
                <div className="my-review-movie">
                  <Link to={`/movie/${review.movie_id}`} className="movie-link">
                    {review.movie?.title || `Фильм #${review.movie_id}`}
                  </Link>
                </div>
                <ReviewCard
                  review={review}
                  onDelete={() => handleDeleteReview(review)}
                  canEdit={true}
                />
              </div>
            ))}
            {totalPages > 1 && (
              <div className="review-list-pagination">
                <button
                  className="pagination-btn"
                  onClick={() => setCurrentPage(currentPage - 1)}
                  disabled={currentPage === 1}
                >
                  Назад
                </button>
                <span className="pagination-info">
                  Страница {currentPage} из {totalPages}
                </span>
                <button
                  className="pagination-btn"
                  onClick={() => setCurrentPage(currentPage + 1)}
                  disabled={currentPage === totalPages}
                >
                  Вперед
                </button>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
};
