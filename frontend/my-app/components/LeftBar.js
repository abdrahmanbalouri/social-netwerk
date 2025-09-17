"use client";

export default function LeftBar({ showSidebar }) {
  return (
    <div className="leftBar">
      <div className="container">
        <div className="menu">
          <div className="user">
            <img
              src="avatar.png"
              alt=""
            />
            <span>user Name</span>
          </div>
          <div className="item">
            <img src="/icone/1.png" alt="" />
            <span>Friends</span>
          </div>
          <div className="item">
            <img src="/icone/2.png" alt="" />
            <span>Groups</span>
          </div>
          <div className="item">
            <img src="/icone/3.png" alt="" />
            <span>Marketplace</span>
          </div>
          <div className="item">
            <img src="/icone/4.png" alt="" />
            <span>Watch</span>
          </div>
          <div className="item">
            <img src="/icone/5.png" alt="" />
            <span>Memories</span>
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
            <img src="/icone/9.png" alt="" />
            <span>Videos</span>
          </div>
          <div className="item">
            <img src="/icone/10.png" alt="" />
            <span>Messages</span>
          </div>
        </div>
        <hr />
        <div className="menu">
          <span>Others</span>
          <div className="item">
            <img src="/icone/11.png" alt="" />
            <span>Fundraiser</span>
          </div>
          <div className="item">
            <img src="/icone/12.png" alt="" />
            <span>Tutorials</span>
          </div>
          <div className="item">
            <img src="/icone/13.png" alt="" />
            <span>Courses</span>
          </div>
        </div>
      </div>
    </div>
  );
}