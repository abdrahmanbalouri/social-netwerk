// components/coment.js
"use client";
import { useState, useEffect, useRef } from "react";
// import { useDarkMode } from '../context/darkMod';
import "../styles/comment.css"
export default function Comment({
  comments,
  isOpen,
  onClose,
  postId,
  postTitle = "",
  onCommentChange,
  lodinggg,
  ongetcomment,
  post
}) {
  //const { darkMode } = useDarkMode();
  const [commentContent, setCommentContent] = useState("");
  const [loading, setLoading] = useState(false);
  const modalRef = useRef(null);
  const [scrollPos, setScrollPos] = useState(0);
  useEffect(() => {

    if (!modalRef.current) return;

    const modal = modalRef.current;

    const reachedBottom = modal.scrollTop + modal.clientHeight >= modal.scrollHeight - 5;
    const previousScrollHeight = modal.scrollHeight;
    async function getcomment() {
      let b = await ongetcomment(post);

      if (b) {
        
  setTimeout(() => {
            const newScrollHeight = modal.scrollHeight;
            const heightIncrease = newScrollHeight - previousScrollHeight;
            
            modal.scrollTop -= heightIncrease - modal.clientHeight ; 
        }, 0); 
      }

    }
    console.log(lodinggg);
    
    if (reachedBottom && !lodinggg) {
      getcomment()
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
      const res = await response.json()

      setCommentContent("");
      if (onCommentChange) {

        onCommentChange(res.comment_id);
        modalRef.current.scrollTop = 0;

      }
    } catch (err) {
      console.error("Error posting comment:", err);
      alert("Failed to post comment. Please try again.");
    } finally {
      setLoading(false);
    }
  }

  if (!isOpen) return null;

  return (
    <div
      className={`modal-overlay ${isOpen ? "is-open" : ""}`}
      onClick={onClose}
    >
      <div
        className="modal-content"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Modal Header */}
        <div className="modal-header">
          <h3 className="modal-title">Comments</h3>
          <button
            className="modal-close"
            aria-label="Close modal"
            onClick={onClose}
          >
            ✕
          </button>
        </div>

        {/* Comments Section */}
        <div className="comments-section">
          {/* Post Title */}

          {/* Comments List */}
          <div
            className="comments-container"
            ref={modalRef}
            onScroll={(e) => setScrollPos(e.target.scrollTop)}

          >

            {comments && comments.length > 0 ? (
              comments.map((comment) => (
                <div key={comment.id} className="comment-item">
                  <div className="comment-header">
                    <span className="comment-author">{comment.author}</span>
                    <span className="comment-date">
                      {new Date(comment.created_at).toLocaleString()}
                    </span>
                  </div>
                  <p className="comment-content">{comment.content}</p>
                </div>
              ))
            ) : (
              <p className="no-comments">No comments yet. Be the first to comment!</p>
            )}
          </div>

          {/* Add Comment Form */}
          <form onSubmit={handlePostComment} className="comment-form">
            <textarea
              className="comment-textarea"
              placeholder="Write a comment..."
              value={commentContent}
              onChange={(e) => setCommentContent(e.target.value)}
              rows={3}
              required
              disabled={loading}
            />
            <button
              type="submit"
              className="post-comment-btn"
              disabled={loading || !commentContent.trim()}
            >
              {loading ? "Posting..." : "Post Comment"}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}