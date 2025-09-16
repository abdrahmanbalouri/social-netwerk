"use client";
import { useEffect, useState } from 'react';
import './profile.css';

export default function Profile() {
  const [userPosts, setUserPosts] = useState([]);
  const [userInfo, setUserInfo] = useState({
    name: 'John Doe',
    friends: 542,
    location: 'New York, NY',
    occupation: 'Software Developer',
    education: 'University of Technology',
    joinedDate: 'Joined January 2020'
  });

  useEffect(() => {
    // Fetch user posts here
    // This is where you'll integrate with your backend
  }, []);

  return (
    <div className="profile-container">
      <div className="profile-header">
        <div className="cover-photo">
          <img src="/default-cover.jpg" alt="Cover" />
        </div>
        
        <div className="profile-info">
          <div className="profile-picture-large">
            <img src="/avatar.png" alt="Profile" />
          </div>
          
          <div className="profile-nav">
            <div>
              <div className="profile-name">
                <h1>{userInfo.name}</h1>
                <div className="friends-count">{userInfo.friends} friends</div>
              </div>
            </div>
            
            <div className="profile-actions">
              <button className="action-button primary-button">
                <i className="fa-solid fa-plus"></i> Add to Story
              </button>
              <button className="action-button">
                <i className="fa-solid fa-pen"></i> Edit Profile
              </button>
            </div>
          </div>
          
          <div className="profile-nav-tabs">
            <button className="active">Posts</button>
            <button>About</button>
            <button>Friends</button>
            <button>Photos</button>
            <button>Videos</button>
          </div>
        </div>
      </div>

      <div className="profile-content">
        <div className="profile-left-column">
          <div className="info-card">
            <h2>Intro</h2>
            <div className="info-list">
              <div className="info-item">
                <i className="fa-solid fa-briefcase"></i>
                {userInfo.occupation}
              </div>
              <div className="info-item">
                <i className="fa-solid fa-graduation-cap"></i>
                {userInfo.education}
              </div>
              <div className="info-item">
                <i className="fa-solid fa-location-dot"></i>
                {userInfo.location}
              </div>
              <div className="info-item">
                <i className="fa-solid fa-clock"></i>
                {userInfo.joinedDate}
              </div>
            </div>
          </div>

          <div className="info-card">
            <h2>Photos</h2>
            {/* Add photo grid here */}
          </div>

          <div className="info-card">
            <h2>Friends</h2>
            {/* Add friends grid here */}
          </div>
        </div>

        <div className="profile-posts">
          <div className="create-post-card">
            <div className="create-post-header">
              <div className="profile-picture">
                <img src="/avatar.png" alt="Profile" />
              </div>
              <input 
                type="text" 
                className="create-post-input"
                placeholder="What's on your mind?"
                readOnly
                onClick={() => {/* Open create post modal */}}
              />
            </div>
          </div>

          {userPosts.map(post => (
            <div key={post.id} className="post">
              <div className="post-header">
                <div className="profile-picture">
                  <img src="/avatar.png" alt="Profile" />
                </div>
                <div>
                  <span className="text-bold">{userInfo.name}</span>
                  <div className="text-muted" style={{ fontSize: '0.85rem' }}>
                    {new Date(post.created_at).toLocaleString()}
                  </div>
                </div>
              </div>
              
              <div className="post-title">{post.title}</div>
              <div className="post-content">{post.content}</div>
              
              {post.image_path && (
                <div className="postimage">
                  <img src={post.image_path} alt="Post content" />
                </div>
              )}
              
              <div className="post-actions">
                <span>Like</span>
                <span>Comment</span>
                <span>Share</span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}