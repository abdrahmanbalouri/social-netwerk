"use client";

import { useRouter } from 'next/navigation';
import { useProfile } from '../context/profile';
import Link from 'next/link';


export default function LeftBar({ showSidebar }) {
  const router = useRouter();
  const { Profile } = useProfile();

  async function handleLogout(e) {
    e?.preventDefault?.();
    try {
      await fetch('/api/logout', { method: 'POST', credentials: 'include' });
    } catch (err) {
      // ignore network errors here; still redirect
      console.error('Logout failed', err);
    }
    // Instead of router.push("/home")
    window.location.href = "/login";
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
              src={Profile?.image ? `/uploads/${Profile.image}` : '/uploads/default.png'}
              alt="user avatar"
            />
            <span>{Profile?.nickname ?? 'user name'}</span>
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
            <Link href={'/watch'}>
              <img src="/icone/4.png" alt="" />
            </Link>
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
            <Link href={'/games'}>
              <img src="/icone/7.png" alt="" />
            </Link>
            <span>Gaming</span>
          </div>
          <div className="item">
            <Link href={'/Gallery'}>
              <img src="/icone/8.png" alt="" />
            </Link>
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