"use client";
import { useState } from 'react';
import Navbar from '../../components/Navbar.js';
import { useDarkMode } from '../../context/darkMod';
import './profile.css';
import LeftBar from '../../components/LeftBar.js';
import RightBar from '../../components/RightBar.js';

export default function Profile() {

  const { darkMode } = useDarkMode();

  const [showSidebar, setShowSidebar] = useState(true);

  async function logout(e) {
    e?.preventDefault?.();
    try {
      const res = await fetch('http://localhost:8080/api/logout', {
        method: 'POST',
        credentials: 'include',
      });
      if (res.ok) router.replace('/login');
    } catch (err) {
      console.error(err);
    }
  }

  return (
    <div className={darkMode ? 'theme-dark' : 'theme-light'}>
      <Navbar
        onLogout={logout}
        onCreatePost={() => setShowModal(true)}
        showSidebar={showSidebar}
        onToggleSidebar={() => setShowSidebar(!showSidebar)}
      />
      <main className="content">
        <LeftBar showSidebar={showSidebar} />

        <div className="profile">
          <div className="images">
            <img
              src="https://images.pexels.com/photos/13440765/pexels-photo-13440765.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=2"
              alt=""
              className="cover"
            />
            <img
              src="avatar.png"
              alt=""
              className="profilePic"
            />
          </div>

          <div className="profileContainer">
            <div className="uInfo">
              <div className="center">
                <span>user name</span>
                <div className="info">
                  <div className="item">
                    <span>about</span>
                  </div>

                </div>
                <button>follow</button>
              </div>
              <div className="right"></div>
            </div>
          </div>
        </div>
        <RightBar />
      </main>
    </div>
  );
}