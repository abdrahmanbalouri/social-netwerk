"use client";
import { useEffect, useState } from 'react';
import Navbar from '../../../components/Navbar.js';
import { useDarkMode } from '../../../context/darkMod.js';
import './profile.css';
import LeftBar from '../../../components/LeftBar.js';
import RightBar from '../../../components/RightBar.js';
import { useParams, useRouter } from 'next/navigation.js';

export default function Profile() {

  const { darkMode } = useDarkMode();

  const [showSidebar, setShowSidebar] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const params = useParams();
  const router = useRouter();
  const userId = Number(params.id);

  const [profile, setProfile] = useState(null);

  async function loadProfile() {
    try {
      const res = await fetch(`http://localhost:8080/api/profile?userId=${params.id}`, {
        method: "GET",
        credentials: "include",
      });
      if (res.ok) {
        const json = await res.json();
        console.log('11',json);
        
        setProfile(json);
      }
    } catch (err) {
      console.error("loadProfile", err);
    }
  }

  useEffect(() => {
    loadProfile();
  }, []);
  
  if (!profile) {
    return (
      <div className={darkMode ? 'theme-dark' : 'theme-light'}>
        <Navbar
          onCreatePost={() => setShowModal(true)}
          onToggleSidebar={() => setShowSidebar(!showSidebar)}
        />
        <main className="content">
          <LeftBar showSidebar={showSidebar} />
          <div className="profile">
            <div className="images">
              <div className="cover" style={{width: '100%', height: 300, background: '#eee'}} />
              <div className="profilePic" style={{width: 200, height: 200, background: '#ccc', borderRadius: '50%', margin: '0 auto', marginTop: -100}} />
            </div>
            <div className="profileContainer">
              <div className="uInfo">
                <div className="center">
                  <span>Loading...</span>
                </div>
              </div>
            </div>
          </div>
          <RightBar />
        </main>
      </div>
    );
  }
  const data = profile;
  return (
    <div className={darkMode ? 'theme-dark' : 'theme-light'}>
      <Navbar
        onCreatePost={() => setShowModal(true)}
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
              src={data.image ? `/uploads/${data.image}` : "avatar.png"}
              alt="profile picture"
              className="profilePic"
            />

          </div>

          <div className="profileContainer">
            <div className="uInfo">
              <div className="center">
                <span>{data.nickname}</span>
                <div className="info">
                  <div className="item">
                    <span>{data.about}</span>
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