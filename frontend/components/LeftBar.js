"use client";

import { useRouter } from 'next/navigation';
import { useProfile } from '../context/profile';
import Link from 'next/link';
import { useWS } from "../context/wsContext.js";
import "../styles/leftbar.css"


export default function LeftBar({ showSidebar }) {
  const router = useRouter();
  const { Profile } = useProfile();
  const { sendMessage } = useWS();
  async function handleLogout(e) {
    e?.preventDefault?.();
    sendMessage({ type: "logout" })
    try {
      await fetch('http://localhost:8080/api/logout', { method: 'POST', credentials: 'include' });
    } catch (err) {
      // ignore network errors here; still redirect
      console.error('Logout failed', err);
    }
    // Instead of router.push("/home")
    window.location.href = "/login";
  }


  return (
    <div className="leftBar" id="leftBar" style={{ display: showSidebar ? 'block' : 'none' }}>
      <div className="menu">
        <div className="user" onClick={() => router.push("/profile/0")}
        >
          <img
            src={Profile?.image ? `/uploads/${Profile.image}` : '/assets/default.png'}
            alt="user avatar"
          />
          <span>{Profile?.first_name ? `${Profile.first_name} ${Profile.last_name}` : 'userName'}</span>
        </div>
        <Link href={`/follow/${Profile?.id}?tab=followers`} >
          <div className="item">

            <img src="/icon/1.png" alt="" />

            <span>followers & following</span>
          </div>
        </Link>
        <Link href='/groups'>
          <div className="item">
            <img src="/icon/2.png" alt="" />
            <span>Groups</span>
          </div>
        </Link>
        <Link href='/watch'>
          <div className="item">
            <img src="/icon/4.png" alt="" />
            <span>Watch</span>
          </div>
        </Link>

      </div>
      <hr />
      <div className="menu">
        <Link href='/Events'>
          <div className="item">
            <img src="/icon/6.png" alt="" />
            <span>Events</span>
          </div>
        </Link>
        <Link href='/games'>
          <div className="item">
            <img src="/icon/7.png" alt="" />
            <span>Gaming</span>
          </div>
        </Link>
        <Link href='/Gallery'>
          <div className="item">
            <img src="/icon/8.png" alt="" />
            <span>Gallery</span>
          </div>
        </Link>
        <Link href="/chat/0">
          <div className="item">
            <img src="/icon/10.png" alt="" />
            <span>Messages</span>
          </div>
        </Link>
      </div>

      <hr />

      <div>
        <button onClick={handleLogout} aria-label="Logout">
          <i className="fa-solid fa-right-from-bracket" />
          <span>Logout</span>
        </button>
      </div>
    </div>
  );
}