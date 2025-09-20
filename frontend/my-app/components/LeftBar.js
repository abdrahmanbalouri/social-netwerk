"use client";

import { useRouter } from 'next/navigation';
import { useProfile } from '../context/profile';


export default function LeftBar({ showSidebar }) {
  const router = useRouter();
  const { profile } = useProfile();

  async function handleLogout(e) {
    e?.preventDefault?.();
    try {
      await fetch('/api/logout', { method: 'POST', credentials: 'include' });
    } catch (err) {
      // ignore network errors here; still redirect
      console.error('Logout failed', err);
    }
    router.replace('/login');
  }

  const logoutStyle = {
    width: '100%',
    height: '40px',
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
    background: 'var(--color-danger, #e53e3e)',
    color: 'white',
    border: 'none',
    padding: '8px 12px',
    borderRadius: '8px',
    cursor: 'pointer',
    fontWeight: 600,
    marginTop: '12px'
  };

  return (
    <div className="leftBar">
      <div className="container">
        <div className="menu">
          <div className="user">
            <img
              src={profile?.image ? `/uploads/${profile.image}` : '/avatar.png'}
              alt="user avatar"
            />
            <span>{profile?.nickname ?? 'user name'}</span>
          </div>
          <div className="item">
            <img src="/icone/1.png" alt="" />
            <span>following</span>
          </div>
          <div className="item">
            <img src="/icone/1.png" alt="" />
            <span>followers</span>
          </div>
          <div className="item">
            <img src="/icone/2.png" alt="" />
            <span>Groups</span>
          </div>
          <div className="item">
            <img src="/icone/4.png" alt="" />
            <span>Watch</span>
          </div>

        </div>
        <hr />
        <div className="menu">
          <span>Your shortcuts</span>
          <div className="item">
            <img src="/icone/6.png" alt="" />
            <span>Events</span>
          </div>
          <div className="item">
            <img src="/icone/7.png" alt="" />
            <span>Gaming</span>
          </div>
          <div className="item">
            <img src="/icone/8.png" alt="" />
            <span>Gallery</span>
          </div>
          <div className="item">
            <img src="/icone/10.png" alt="" />
            <span>Messages</span>
          </div>
        </div>
        <hr />

        <div style={{ display: 'flex', justifyContent: 'center' }}>
          <button onClick={handleLogout} style={logoutStyle} aria-label="Logout">
            <i className="fa-solid fa-right-from-bracket" />
            <span>Logout</span>
          </button>
        </div>
      </div>
    </div>
  );
}