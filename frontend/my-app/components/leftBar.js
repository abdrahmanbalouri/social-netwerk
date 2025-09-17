"use client";

export default function LeftBar({ showSidebar }) {
  return (
    <aside className={`sidebar ${showSidebar ? 'show' : 'hide'}`}>
      <ul>
        <li><a><i className="fa-solid fa-house"></i>Home</a></li>
        <li><a><i className="fa-solid fa-user-group"></i>Friends</a></li>
        <li><a><i className="fa-solid fa-users"></i>Groups</a></li>
        <li><a><i className="fa-solid fa-gamepad"></i>games</a></li>
        <li><a><i className="fa-solid fa-video"></i>reel</a></li>
        <li><a><i className="fa-solid fa-clock"></i>Memories</a></li>
      </ul>
    </aside>
  );
}