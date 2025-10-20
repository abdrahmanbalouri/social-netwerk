"use client"
import { useState, useEffect, useRef } from "react"
import "../styles/comment.css"

export default function Comment({ comments, isOpen, onClose, postId, onCommentChange, lodinggg, ongetcomment, post }) {
  const [commentContent, setCommentContent] = useState("")
  const [selectedFile, setSelectedFile] = useState(null)
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

  useEffect(()=>{

  console.log(comments);
  
  },[comments])

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
    if (!commentContent.trim() && !selectedFile) return

    try {
      setLoading(true)
      const formData = new FormData()
      formData.append("postId", postId)
      formData.append("content", commentContent)
      if (selectedFile) {
        formData.append("media", selectedFile)
      }

      const response = await fetch("http://localhost:8080/api/createcomment", {
        method: "POST",
        credentials: "include",
        body: formData,
      })

      if (!response.ok) {
        throw new Error("Failed to post comment")
      }
      

      const res = await response.json()
      setCommentContent("")
      setSelectedFile(null)

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
                  src={`../${post.image_path}`}
                  alt="Post"
                  className="cm-img"
                  onLoad={(e) => e.target.classList.add("cm-loaded")}
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

                        {comment.media_path && (
                          comment.media_path.endsWith(".mp4") || comment.media_path.endsWith(".webm") ? (
                            <video controls className="cm-media">
                              <source src={`../${comment.media_path}`} type="video/mp4" />
                              Your browser does not support the video tag.
                            </video>
                          ) : (
                            <img src={`../${comment.media_path}`} alt="Comment Media" className="cm-media" />
                          )
                        )}
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
              <form onSubmit={handlePostComment} className="cm-form" encType="multipart/form-data">
                <div className="cm-input-wrap">
                  <textarea
                    className="cm-input"
                    placeholder="Add a comment..."
                    value={commentContent}
                    onChange={(e) => setCommentContent(e.target.value)}
                    rows={1}
                    required={!selectedFile}
                    disabled={loading}
                  />
                  <input
                    type="file"
                    accept="image/*,video/*"
                    onChange={(e) => setSelectedFile(e.target.files[0])}
                    className="cm-file"
                  />
                  <button
                    type="submit"
                    className="cm-btn"
                    disabled={loading || (!commentContent.trim() && !selectedFile)}
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
  )
}
