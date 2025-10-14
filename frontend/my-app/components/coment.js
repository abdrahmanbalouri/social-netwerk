// components/comment.js
"use client";
import { useState, useEffect, useRef } from "react";
import "../styles/comment.css";

export default function Comment({
  comments,
  isOpen,
  onClose,
  postId,
  onCommentChange,
  lodinggg,
  ongetcomment,
  post
}) {
  const [commentContent, setCommentContent] = useState("");
  const [loading, setLoading] = useState(false);
  const modalRef = useRef(null);
  const commentsContainerRef = useRef(null);
  const [scrollPos, setScrollPos] = useState(0);
  const commentRefs = useRef({});
  
  function scrollToComment(commentId) {
    const el = commentRefs.current[commentId];
    if (el) {
      el.scrollIntoView({ behavior: "smooth", block: "start" });
    }
  }
  
  useEffect(() => {
    if (!commentsContainerRef.current) return;

    const commentsContainer = commentsContainerRef.current;
    const reachedBottom = commentsContainer.scrollTop + commentsContainer.clientHeight >= commentsContainer.scrollHeight - 5;

    async function getcomment() {
      let r = await ongetcomment(post);
      if (r) {
        scrollToComment(r);
      }
    }

    if (reachedBottom && !lodinggg) {
      getcomment();
    }
  }, [scrollPos]);

  async function handlePostComment(e) {
    e.preventDefault();

    if (!commentContent.trim()) return;

    try {
      setLoading(true);
      const response = await fetch("http://localhost:8080/api/createcomment", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          postId: postId,
          content: commentContent,
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to post comment");
      }
      const res = await response.json();

      setCommentContent("");
      if (onCommentChange) {
        onCommentChange(res.comment_id);
        commentsContainerRef.current.scrollTop = 0;
      }
    } catch (err) {
      console.error("Error posting comment:", err);
      alert("Failed to post comment. Please try again.");
    } finally {
      setLoading(false);
    }
  }

  if (!isOpen) return null;

  // Determine if we should show the image section
  const showImageSection = !!post.image_path;

  return (
    <div
      className={`modal-overlay ${isOpen ? "is-open" : ""}`}
      onClick={onClose}
    >
      <div
        className="modal-content"
        onClick={(e) => e.stopPropagation()}
        ref={modalRef}
      >
        {/* Modal Header */}
        <div className="modal-header">
          <button
            className="modal-close"
            aria-label="Close modal"
            onClick={onClose}
          >
            ✕
          </button>
          <h3 className="modal-title">Comments</h3>
        </div>

        <div className={`modal-body ${showImageSection ? "" : "no-image-mode"}`}>
          {/* Post Section - Left Side (Conditional Render) */}
          {showImageSection && (
            <div className="post-section">
              <div className="image-container">
                <img 
                  src={post.image_path} 
                  alt="Post" 
                  className="post-image"
                  onLoad={(e) => {
                    e.target.classList.add('image-loaded');
                  }}
                />
              </div>
            </div>
          )}

          {/* Comments Section - Right Side (Always full width if no image) */}
          <div className="comments-section">
            {/* Post Header in Comments */}
            <div className="post-header-comments">
              <div className="user-info">
                <div className="user-avatar">
                  {post.author?.charAt(0) || "U"}
                </div>
                <div className="user-details">
                  <span className="username">{post.author || "Unknown"}</span>
                </div>
              </div>
              <div className="post-title">
                <span className="caption-text">{post.title}</span>
              </div>
              {post.content && (
                <div className="post-caption">
                  {post.content}
                </div>
              )}
            </div>

            {/* Comments List */}
            <div
              className="comments-container"
              ref={commentsContainerRef}
              onScroll={(e) => setScrollPos(e.target.scrollTop)}
            >
              {comments && comments.length > 0 ? (
                <div className="comments-list">
                  {comments.map((comment) => (
                    <div
                      key={comment.id}
                      id={`comment-${comment.id}`}
                      className="comment"
                      ref={(el) => (commentRefs.current[comment.id] = el)}
                    >
                      <div className="comment-avatar">
                        {comment.author?.charAt(0) || "U"}
                      </div>
                      <div className="comment-content">
                        <div className="comment-header">
                          <span className="comment-author">{comment.author}</span>
                          <span className="comment-time">
                            {new Date(comment.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                          </span>
                        </div>
                        <p className="comment-text">{comment.content}</p>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="no-comments">
                  <div className="no-comments-icon">💬</div>
                  <p className="no-comments-title">No comments yet</p>
                  <p className="no-comments-subtitle">Start the conversation</p>
                </div>
              )}
            </div>

            {/* Comment Form */}
            <div className="comment-form-container">
              <form onSubmit={handlePostComment} className="comment-form">
                <div className="input-container">
                  <textarea
                    className="comment-input"
                    placeholder="Add a comment..."
                    value={commentContent}
                    onChange={(e) => setCommentContent(e.target.value)}
                    rows={1}
                    required
                    disabled={loading}
                  />
                  <button
                    type="submit"
                    className="post-comment-btn-header"
                    disabled={loading || !commentContent.trim()}
                  >
                    {loading ? "..." : "Post"}
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}