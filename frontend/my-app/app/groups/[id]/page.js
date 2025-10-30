"use client";

import "../../../styles/events.css"
import Navbar from "../../../components/Navbar.js";
import { GroupPostChat } from "../../../components/groupPostCat.js";
import Post from "../../../components/Post.js";
import { useEffect, useState } from "react";
import "../../../styles/groupstyle.css";
import { useParams } from "next/navigation";
import { PostCreationTrigger } from "../../../components/cretaePostGroup.js";
import LeftBar from "../../../components/LeftBar.js";
import RightBarGroup from "../../../components/RightBarGroups.js";
import { useDarkMode } from "../../../context/darkMod.js";
import { FileText } from "lucide-react";
// import {CreatePost} from ""

export default function () {
  const { darkMode } = useDarkMode();

  return (
    <div className={darkMode ? "theme-dark" : "theme-light"}>
      <Navbar />
      <main className="content">
        <LeftBar showSidebar={true} />
        {/* <AllPosts /> */}
        <GroupPostChat />
        {/* <RightBar /> */}
        <RightBarGroup onClick={sendRequest} />
      </main>
      {/* <AllPosts /> */}
      {/* <CreatePost /> */}
    </div>
  );
}

function sendRequest(invitedUserID, grpID) {
  console.log("invited user id is : ", invitedUserID);
  console.log("group id is : ", grpID);

  fetch(`http://localhost:8080/group/invitation/${grpID}`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      invitedUsers: [invitedUserID],
    }),
  })
    .then(async (res) => {
      const temp = await res.json();
      console.log("temp is :", temp);
    })
    .catch((error) => {
      console.log("error sending invitation to the user :", error);
    });
}

export function AllPosts() {
  const [posts, setPost] = useState([]);
  const [loading, setLoading] = useState(true);
  const params = useParams();
  const grpID = params.id;

  useEffect(() => {
    if (!grpID) {
      setLoading(false);
      return;
    }

    fetch(`http://localhost:8080/group/fetchPosts/${grpID}`, {
      method: "GET",
      credentials: "include",
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch posts");
        return res.json();
      })
      .then((data) => {
        setPost(data || []);
        setLoading(false);
      })
      .catch((error) => {
        console.error("Failed to fetch posts:", error);
        setLoading(false);
      });
  }, [grpID]);
  // if (!posts) {
  //     return (
  //         <div>
  //             <PostCreationTrigger setPost = {setPost}/>
  //             <div>There is no post yeeeeeet.</div>
  //         </div>
  //     );
  // }

  console.log("posts are :", posts);
  return (
    <div>
      <PostCreationTrigger setPost={setPost} />
      <div className="content-area">
        {posts.length === 0 ? (
          <div className="group-empty-state">
            <FileText />
            <p className="group-empty-state-text">
              No posts yet. Be the first to share something!
            </p>
          </div>
        ) : (
          posts.map((post) => (
            <Post
              key={post.id}
              post={post}
              onGetComments={GetComments}
              ondolike={AddLike}
            />
          ))
        )}
      </div>
    </div>
  );
}
export function LastPost() {
    const [posts, setPosts] = useState([]);
    const [loading, setLoading] = useState(true);
    const params = useParams();
    const grpID = params.id;

    useEffect(() => {
        if (!grpID) {
            setLoading(false);
            return;
        }

        fetch(`http://localhost:8080/group/fetchPost/${grpID}`, {
            method: 'GET',
            credentials: 'include',
        })
            .then(res => {
                if (!res.ok) throw new Error('Failed to fetch last post');
                return res.json();
            })
            .then(data => {
                setPosts(data);
                setLoading(false);
            })
            .catch(error => {
                console.error("Failed to fetch last post:", error);
                setLoading(false);
            });
    }, [grpID]);

    if (loading) {
        return <div>Loading...</div>;
    }

    if (!posts || posts.length === 0) {
        return (
            <div>
                <PostCreationTrigger setPosts={setPosts} />
                <div>There is no post yet.</div>
            </div>
        );
    }

    console.log("posts are:", posts);

    return (
        <div>
            <PostCreationTrigger setPosts={setPosts} />
            <div className="content-area">
                {posts.map((post) => (
                    <Post
                        key={post.id}
                        post={post}
                        onGetComments={GetComments}
                        ondolike={AddLike}
                    />
                ))}
            </div>
        </div>
    );
}

function GetComments(post) {
  // setSelectedPost({
  //     id: post.id,
  //     title: post.title || post.post_title || "Post"
  // });
}

function AddLike() {}

export async function CreatePost(groupId, formData) {
  const data = new FormData();
  data.append("postData", JSON.stringify(formData));

  const res = await fetch(`http://localhost:8080/group/addPost/${groupId}`, {
    method: "POST",
    credentials: "include",
    body: data,
  });

  const text = await res.text();

  if (!res.ok) throw new Error("Failed to create a post htmlFor groups");

  return JSON.parse(text);
}


// events 




export  function Events() {
  return(
    <>
   <div className="events-container">
  <div className="create-event-card">
    <h2 className="card-title">Create Event</h2>
    <form className="event-form">
      <div className="form-group">
        <label htmlFor="event-title">Title</label>
        <input type="text" id="event-title" placeholder="Event title..." />
      </div>
      
      <div className="form-group">
        <label htmlFor="event-description">Description</label>
        <textarea id="event-description" rows="3" placeholder="Event description..."></textarea>
      </div>
      
      <div className="form-group">
        <label htmlFor="event-datetime">Day/Time</label>
        <input type="datetime-local" id="event-datetime" />
      </div>
      
      <button type="submit" className="btn-create">Create Event</button>
    </form>
  </div>

  <div className="events-list">
    <div className="event-card">
      <div className="event-header">
        <h3 className="event-title">Team Meeting</h3>
        <span className="event-datetime">Nov 15, 2024 - 14:00</span>
      </div>
      <p className="event-description">
        Monthly team sync to discuss project progress and upcoming milestones.
      </p>
      <div className="event-actions">
        <button className="btn-going">Going</button>
        <button className="btn-not-going">Not Going</button>
      </div>
      <div className="event-stats">
        <span className="stat-going">12 going</span>
        <span className="stat-not-going">3 not going</span>
      </div>
    </div>

    <div className="event-card">
      <div className="event-header">
        <h3 className="event-title">Code Review Session</h3>
        <span className="event-datetime">Nov 20, 2024 - 10:00</span>
      </div>
      <p className="event-description">
        Review and discuss recent code changes with the development team.
      </p>
      <div className="event-actions">
        <button className="btn-going">Going</button>
        <button className="btn-not-going">Not Going</button>
      </div>
      <div className="event-stats">
        <span className="stat-going">8 going</span>
        <span className="stat-not-going">2 not going</span>
      </div>
    </div>
  </div>
</div>
    </>
  )
}