"use client";

import { useRouter } from 'next/navigation';
import { useProfile } from '../context/profile';
import Link from 'next/link';
import "../styles/leftbar.css"


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
console.log(Profile);

  return (
    <div className="leftBar">
        <div className="menu">
          <div className="user">
            <img
              src={Profile?.image ? `/uploads/${Profile.image}` : '/uploads/default.png'}
              alt="user avatar"
            />
            <span>{Profile?.nickname ?? 'user name'}</span>
          </div>
          <Link href={`/follow/${Profile?.id}?tab=following`}  >  
          <div className="item">
            <img src="/icone/1.png" alt="" />
            <span>following</span>
          </div>
            </Link>
          <Link href={`/follow/${Profile?.id}?tab=followers`} >  
          <div className="item">

            <img src="/icone/1.png" alt="" />

            <span>followers</span>
          </div>
            </Link>
          <div className="item">
            <img src="/icone/2.png" alt="" />
            <span>Groups</span>
          </div>
            <Link href={'/watch'}>
          <div className="item">
              <img src="/icone/4.png" alt="" />
            <span>Watch</span>
          </div>
            </Link>

        </div>
        <hr />
        <div className="menu">
          <span>Your shortcuts</span>
          <div className="item">
            <img src="/icone/6.png" alt="" />
            <span>Events</span>
          </div>
            <Link href={'/games'}>
          <div className="item">
              <img src="/icone/7.png" alt="" />
            <span>Gaming</span>
          </div>
            </Link>
            <Link href={'/Gallery'}>
          <div className="item">
              <img src="/icone/8.png" alt="" />
            <span>Gallery</span>
          </div>
            </Link>
          <div className="item">
            <img src="/icone/10.png" alt="" />
            <span>Messages</span>
          </div>
        </div>

        <div style={{ display: 'flex', justifyContent: 'center' }}>
          <button onClick={handleLogout} style={logoutStyle} aria-label="Logout">
            <i className="fa-solid fa-right-from-bracket" />
            <span>Logout</span>
          </button>
        </div>
    </div>
  );
}