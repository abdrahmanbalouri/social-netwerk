"use client"
import { useState, useEffect, useRef } from "react"
import "../styles/comment.css"

export default function Comment({ comments, isOpen, onClose, postId, onCommentChange, lodinggg, ongetcomment, post }) {
  const [commentContent, setCommentContent] = useState("")
  const [loading, setLoading] = useState(false)
  const modalRef = useRef(null)
  const commentsContainerRef = useRef(null)
  const [scrollPos, setScrollPos] = useState(0)
  const commentRefs = useRef({})

  function scrollToComment(commentId) {
    const el = commentRefs.current[commentId]
    if (el) {
      el.scrollIntoView({ behavior: "smooth", block: "start" })
    }
  }

  useEffect(() => {
    if (!commentsContainerRef.current) return

    const commentsContainer = commentsContainerRef.current
    const reachedBottom =
      commentsContainer.scrollTop + commentsContainer.clientHeight >= commentsContainer.scrollHeight - 5

    async function getcomment() {
      const r = await ongetcomment(post)
      if (r) {
        scrollToComment(r)
      }
    }

    if (reachedBottom && !lodinggg) {
      getcomment()
    }
  }, [scrollPos])

  async function handlePostComment(e) {
    e.preventDefault()

    if (!commentContent.trim()) return

    try {
      setLoading(true)
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
      })

      if (!response.ok) {
        throw new Error("Failed to post comment")
      }
      const res = await response.json()

      setCommentContent("")
      if (onCommentChange) {
        onCommentChange(res.comment_id)
        commentsContainerRef.current.scrollTop = 0
      }
    } catch (err) {
      console.error("Error posting comment:", err)
      alert("Failed to post comment. Please try again.")
    } finally {
      setLoading(false)
    }
  }

  if (!isOpen) return null

  const showImageSection = !!post.image_path

  return (
    <div className={`cm-overlay ${isOpen ? "cm-active" : ""}`} onClick={onClose}>
      <div className="cm-modal" onClick={(e) => e.stopPropagation()} ref={modalRef}>
        <div className="cm-header">
          <button className="cm-close" aria-label="Close modal" onClick={onClose}>
            <span className="icon-close"></span>
          </button>
          <h3 className="cm-heading">Comments</h3>
        </div>

        <div className={`cm-body ${showImageSection ? "" : "cm-full"}`}>
          {showImageSection && (
            <div className="cm-post">
              <div className="cm-img-wrap">
                <img
                  src={post.image_path || "/placeholder.svg"}
                  alt="Post"
                  className="cm-img"
                  onLoad={(e) => {
                    e.target.classList.add("cm-loaded")
                  }}
                />
              </div>
            </div>
          )}

          <div className="cm-content">
            <div className="cm-author">
              <div className="cm-user">
                <div className="cm-avatar">{post.author?.charAt(0) || "U"}</div>
                <div className="cm-info">
                  <span className="cm-name">{post.author || "Unknown"}</span>
                </div>
              </div>
              <div className="cm-caption">
                <span className="cm-title">{post.title}</span>
              </div>
              {post.content && <div className="cm-desc">{post.content}</div>}
            </div>

            <div
              className="cm-list-wrap"
              ref={commentsContainerRef}
              onScroll={(e) => setScrollPos(e.target.scrollTop)}
            >
              {comments && comments.length > 0 ? (
                <div className="cm-list">
                  {comments.map((comment) => (
                    <div
                      key={comment.id}
                      id={`comment-${comment.id}`}
                      className="cm-item"
                      ref={(el) => (commentRefs.current[comment.id] = el)}
                    >
                      <div className="cm-avatar cm-avatar-sm">{comment.author?.charAt(0) || "U"}</div>
                      <div className="cm-text">
                        <div className="cm-meta">
                          <span className="cm-name">{comment.author}</span>
                          <span className="cm-time">
                            {new Date(comment.created_at).toLocaleTimeString([], {
                              hour: "2-digit",
                              minute: "2-digit",
                            })}
                          </span>
                        </div>
                        <p className="cm-msg">{comment.content}</p>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="cm-empty">
                  <div className="icon-chat"></div>
                  <p className="cm-empty-title">No comments yet</p>
                  <p className="cm-empty-text">Start the conversation</p>
                </div>
              )}
            </div>

            <div className="cm-form-wrap">
              <form onSubmit={handlePostComment} className="cm-form">
                <div className="cm-input-wrap">
                  <textarea
                    className="cm-input"
                    placeholder="Add a comment..."
                    value={commentContent}
                    onChange={(e) => setCommentContent(e.target.value)}
                    rows={1}
                    required
                    disabled={loading}
                  />
                  <button type="submit" className="cm-btn" disabled={loading || !commentContent.trim()}>
                    {loading ? "..." : "Post"}
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
