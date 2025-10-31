"use client";

import { useState, useEffect, useRef } from "react";
import { useProfile } from "../context/profile";
import "../styles/stories.css";

const Stories = () => {
  const { Profile } = useProfile();

  const [stories, setStories] = useState([]);
  const [showModal, setShowModal] = useState(false);
  const [newStory, setNewStory] = useState({
    content: "",
    imageFile: null,
    bgColor: "#000000",
  });
  const [viewStoryIndex, setViewStoryIndex] = useState(null);
  const [timeLeft, setTimeLeft] = useState(6);
  const intervalId = useRef(null);

  // Toast state
  const [toast, setToast] = useState(null);

  const showToast = (message, type = "error", duration = 3000) => {
    setToast({ message, type });
    setTimeout(() => {
      setToast(null);
    }, duration);
  };

  const fetchStories = async () => {
    try {
      const res = await fetch("http://localhost:8080/api/Getstories", {
        credentials: "include",
      });

      const response = await res.json()
      if (response.error) {
        //showToast(response.error);
        return
      }
       
      setStories(Array.isArray(response) ? response : []);
    } catch (err) {    
      //showToast("Failed to load stories");
      setStories([]);
    }
  };

  useEffect(() => {
    fetchStories();
    const interval = setInterval(fetchStories, 10000);
    return () => clearInterval(interval);
  }, []);

  const groupStoriesByUser = () => {
    const grouped = {};
    stories.forEach(s => {
      if (!grouped[s.user_id]) {
        grouped[s.user_id] = {
          user: {
            id: s.user_id,
            first_name: s.first_name,
            last_name: s.last_name,
            profile: s.profile,
          },
          stories: [],
        };
      }
      grouped[s.user_id].stories.push(s);
    });
    
    return Object.values(grouped);
  };

  const groupedStories = groupStoriesByUser();
  

const activeGrouped = groupedStories
  .map(group => ({
    ...group,
    stories: group.stories.filter(
      s => !s.expires_at || new Date(s.expires_at) > new Date()
    ),
  }))
  .filter(group => group.stories.length > 0)
  .sort((a, b) => {
    if (!Profile?.id) return 0;
    const aIsMine = a.user?.id === Profile.id;
    const bIsMine = b.user?.id === Profile.id;

    if (aIsMine && !bIsMine) return -1; 
    if (!aIsMine && bIsMine) return 1;  
    return 0;
  });


  const allGroups = [...activeGrouped];


  const myGroup = allGroups.find(g => g.user.id === Profile?.id);
  const hasMyStories = myGroup && myGroup.stories.length > 0;

  const otherGroups = allGroups.filter(g => g.user.id !== Profile?.id);

  const handleCreateStory = async () => {
    try {
      if (!newStory.content && !newStory.imageFile) {
        showToast("Add text or an image");
        return;
      }

      const form = new FormData();
      form.append("content", newStory.content);
      form.append("bg_color", newStory.bgColor);
      if (newStory.imageFile) {
        form.append("image", newStory.imageFile);
      }

      const res = await fetch("http://localhost:8080/api/Createstories", {
        method: "POST",
        credentials: "include",
        body: form,
      });
      const response = await res.json()

      if (response.error) {
        showToast(response.error);
        return

      }

      setNewStory({ content: "", imageFile: null, bgColor: "#000000" });
      setShowModal(false);
      showToast("Story created successfully!", "success");
      await fetchStories();
    } catch (err) {
      showToast(err.message || "Failed to create story");
    }
  };

  const openViewer = (group, storyIndex = 0) => {
    const groupIndex = allGroups.findIndex(g => g.user.id === group.user.id);
    
    
    if (groupIndex !== -1) {
      setViewStoryIndex({ group: groupIndex, story: storyIndex });
      setTimeLeft(6);
    }
  };

  const closeViewer = () => {
    setViewStoryIndex(null);
    setTimeLeft(6);
    if (intervalId.current) {
      clearInterval(intervalId.current);
      intervalId.current = null;
    }
  };

  const startCountdown = () => {
    if (intervalId.current) {
      clearInterval(intervalId.current);
      intervalId.current = null;
    }

    intervalId.current = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev <= 1) {
          goToNextStory();
          return 6;
        }
        return prev - 1;
      });
    }, 1000);
  };

  const goToNextStory = () => {
    if (!viewStoryIndex) return;

    const { group: gIdx, story: sIdx } = viewStoryIndex;
    const currentGroup = allGroups[gIdx];

    if (sIdx + 1 < currentGroup.stories.length) {
      setViewStoryIndex({ group: gIdx, story: sIdx + 1 });
      setTimeLeft(6);
    } else if (gIdx + 1 < allGroups.length) {
      setViewStoryIndex({ group: gIdx + 1, story: 0 });
      setTimeLeft(6);
    } else {
      closeViewer();
    }
  };

  const goToPreviousStory = () => {
    if (!viewStoryIndex) return;

    const { group: gIdx, story: sIdx } = viewStoryIndex;

    if (sIdx > 0) {
      setViewStoryIndex({ group: gIdx, story: sIdx - 1 });
      setTimeLeft(6);
    } else if (gIdx > 0) {
      const prevGroup = allGroups[gIdx - 1];
      setViewStoryIndex({
        group: gIdx - 1,
        story: prevGroup.stories.length - 1
      });
      setTimeLeft(6);
    }
  };

  const handlePrevious = (e) => {
    e.stopPropagation();
    goToPreviousStory();
  };

  const handleNext = (e) => {
    e.stopPropagation();
    goToNextStory();
  };

  useEffect(() => {
    if (!viewStoryIndex) {
      if (intervalId.current) {
        clearInterval(intervalId.current);
        intervalId.current = null;
      }
      return;
    }

    setTimeLeft(6);
    startCountdown();

    return () => {
      if (intervalId.current) {
        clearInterval(intervalId.current);
        intervalId.current = null;
      }
    };
  }, [viewStoryIndex]);

  const currentGroupIdx = viewStoryIndex?.group ?? null;
  const currentStoryIdx = viewStoryIndex?.story ?? null;
  const currentGroup = currentGroupIdx !== null ? allGroups[currentGroupIdx] : null;
  const currentStory = currentGroup && currentStoryIdx !== null ? currentGroup.stories[currentStoryIdx] : null;

  const getCurrentGroupSegments = () => {
    return currentGroup ? currentGroup.stories.length : 0;
  };

  const totalSegments = getCurrentGroupSegments();
  const currentSegment = currentStoryIdx;

  const canGoNext = () => {
    if (!viewStoryIndex) return false;
    const { group: gIdx, story: sIdx } = viewStoryIndex;
    const currentGroup = allGroups[gIdx];
    return (sIdx + 1 < currentGroup.stories.length) || (gIdx + 1 < allGroups.length);
  };

  const canGoPrev = () => {
    if (!viewStoryIndex) return false;
    const { group: gIdx, story: sIdx } = viewStoryIndex;
    return (sIdx > 0) || (gIdx > 0);
  };

  return (
    <div className="story-bar">

      {/* Toast Notification */}
      {toast && (
        <div className={`toast ${toast.type}`}>
          <span>{toast.message}</span>
          <button onClick={() => setToast(null)} className="toast-close">×</button>
        </div>
      )}

      {/* "Your Story" */}
      {Profile && (
        <div
          className={`story-item my-story ${hasMyStories ? 'has-story' : ''}`}
          onClick={() => {
            if (hasMyStories) {
              openViewer(myGroup, 0);
            }
          }}
        >
          <div className="story-image-wrapper">
            <img
              src={Profile?.image ? `/uploads/${Profile.image}` : "/assets/default.png"}
              alt="Your avatar"
              className="profile-img"
            />
          </div>
          <span className="story-label">Your Story</span>

          <button
            className="story-add-btn add-more"
            onClick={(e) => {
              e.stopPropagation();
              setShowModal(true);
            }}
          >
            +
          </button>
        </div>
      )}

      {/* Other Users' Stories */}
      {otherGroups.map((group) => (
        <div
          key={group.user.id}
          className="story-item"
          onClick={() => openViewer(group, 0)}
        >
          <div className="story-image-wrapper">
            <img
              src={group?.user?.profile ? `/uploads/${group.user.profile}` : "/assets/default.png"}
              alt={group.user.first_name}
              className="profile-img"
            />
          </div>
          <span className="story-label">
            {group.user.first_name} {group.user.last_name}
          </span>
        </div>
      ))}

      {/* Create Modal */}
      {showModal && (
        <div className="story-modal-overlay" onClick={() => setShowModal(false)}>
          <div className="story-modal" onClick={e => e.stopPropagation()}>
            <h2>Create Story</h2>
            <div className="create-cover-preview">
              {newStory.imageFile ? (
                <img src={URL.createObjectURL(newStory.imageFile)} alt="Preview" />
              ) : newStory.content ? (
                <div style={{
                  backgroundColor: newStory.bgColor,
                  width: '100%',
                  height: '100%',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  color: '#fff',
                  fontSize: '20px',
                  fontWeight: '600'
                }}>
                  {newStory.content}
                </div>
              ) : (
                <div className="placeholder-cover">Choose Image or Add Text</div>
              )}
            </div>
            <input
              type="text"
              placeholder="Your text (optional)"
              value={newStory.content}
              onChange={e => setNewStory({ ...newStory, content: e.target.value })}
            />
            <input
              type="color"
              value={newStory.bgColor}
              onChange={e => setNewStory({ ...newStory, bgColor: e.target.value })}
            />
            <input
              type="file"
              accept="image/*"
              onChange={e => {
                const file = e.target.files?.[0];
                if (file) setNewStory({ ...newStory, imageFile: file });
              }}
            />
            <div className="modal-buttons">
              <button onClick={handleCreateStory}>Post Story</button>
              <button onClick={() => {
                setShowModal(false);
                setNewStory({ content: "", imageFile: null, bgColor: "#000000" });
              }}>Cancel</button>
            </div>
          </div>
        </div>
      )}

      {/* Story Viewer */}
      {currentStory && (
        <div className="story-view-overlay" onClick={closeViewer}>
          <div className="story-view-container" onClick={e => e.stopPropagation()}>
            {/* Progress Bars */}
            <div className="stories-progress-container">
              {Array.from({ length: totalSegments }).map((_, segIdx) => {
                const isCompleted = segIdx < currentSegment;
                const isActive = segIdx === currentSegment;
                const fillWidth = isCompleted
                  ? "100%"
                  : isActive
                    ? `${((6 - timeLeft) / 6) * 100}%`
                    : "0%";

                return (
                  <div
                    key={`${currentGroupIdx}-${segIdx}`}
                    className={`story-segment-bar ${isCompleted ? "completed" : ""} ${isActive ? "active" : ""}`}
                  >
                    <div
                      className="story-segment-fill"
                      style={{
                        width: fillWidth,
                        transition: isActive ? "width 1s linear" : "none",
                      }}
                    />
                  </div>
                );
              })}
            </div>

            {/* User Info */}
            <div className="story-user-info">
              <img
                src={currentGroup.user.profile ? `/uploads/${currentGroup.user.profile}` : "/avatar.png"}
                alt={currentGroup.user.first_name}
                className="story-user-avatar"
              />
              <span className="story-user-name">
                {currentGroup.user.first_name} {currentGroup.user.last_name}
              </span>
              <span className="story-time">
                {new Date(currentStory.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
              </span>
            </div>

            {/* Timer */}
            <div className="story-timer-text">{timeLeft}s</div>

            {/* Story Content */}
            {currentStory.image_url ? (
              <div className="story-content-container">
                <img
                  src={currentStory.image_url}
                  alt="Story"
                  className="story-centered-image"
                />
              </div>
            ) : (
              <div className="story-content-container">
                <div
                  className="story-text-content"
                  style={{
                    backgroundColor: currentStory.bg_color || '#000',
                    color: currentStory.text_color || '#fff'
                  }}
                >
                  {currentStory.content}
                </div>
              </div>
            )}

            {/* Navigation */}
            {canGoPrev() && (
              <button className="story-nav-btn story-prev-btn" onClick={handlePrevious}>
                ‹
              </button>
            )}
            {canGoNext() && (
              <button className="story-nav-btn story-next-btn" onClick={handleNext}>
                ›
              </button>
            )}

            <div className="tap-left" onClick={handlePrevious} />
            <div className="tap-right" onClick={handleNext} />

            <button className="story-close-btn" onClick={closeViewer}>×</button>
          </div>
        </div>
      )}
    </div>
  );
};

export default Stories;