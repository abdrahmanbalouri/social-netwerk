"use client";

import "../../../styles/events.css"
import Navbar from "../../../components/Navbar.js";
import { GroupPostChat } from "../../../components/groupPostCat.js";
import Post from "../../../components/Post.js";
import { use, useEffect, useState } from "react";
import "../../../styles/groupstyle.css";
import { useParams } from "next/navigation";
import { PostCreationTrigger } from "../../../components/cretaePostGroup.js";
import LeftBar from "../../../components/LeftBar.js";
import RightBarGroup from "../../../components/RightBarGroups.js";
import { useDarkMode } from "../../../context/darkMod.js";
import EventCard from "../../../components/EventCard.js";
import { FileText, MessageCircle } from "lucide-react";
import { useWS } from "../../../context/wsContext.js";
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

function AddLike() { }

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


export function Events() {
  const [ShowEventForm, SetShowEventForm] = useState(true)
  const params = useParams();

  const grpID = params.id;

  const [EventList, setEventList] = useState([])

  useEffect(() => {
    fetchEvents();
  }, []);

  async function fetchEvents() {
    try {
      const response = await fetch(`http://localhost:8080/api/getEvents/${grpID}`, {
        method: 'GET',
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error('Failed to fetch events');
      }

      const data = await response.json();

      console.log('Fetched events:', data);
      setEventList(data);
    } catch (error) {
      console.error('Error fetching events:', error);
    }
  }

  function showEvent() {
    SetShowEventForm(prev => !prev)
  }



  async function goingEvent(status, eventID) {

    if ((!status || !eventID) || (status !== "going" && status !== "notGoing")) {
      console.error('Status and Event ID are required');
      return;
    }
    try {
      const response = await fetch(`http://localhost:8080/api/event/action/${grpID}`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ status, eventID }),
      });

      if (!response.ok) {
        throw new Error('Failed to RSVP to event');
      }

      await fetchEvents();
    } catch (error) {
      console.error('Error RSVPing to event:', error);
    }
  }


  return (
    <>
      {ShowEventForm ? (
        <div>
          <button onClick={showEvent} className="Showbutton">
            create event
          </button>
        </div>
      ) : (
        <div className="privacy-overlay">

          <div className="backLayer" onClick={showEvent}></div>
          <EventForm ShowEventForm={ShowEventForm} closeForm={showEvent}
            fetchEvents={fetchEvents}
          />
        </div>
      )}

      {EventList?.length === 0 ? (
        <div>No events yet. Be the first to create one!</div>
      ) : (
        EventList.map((ev) => (
          <EventCard key={ev.id} ev={ev} goingEvent={goingEvent} />
        ))
      )}





    </>
  )
}



export function EventForm({ closeForm, fetchEvents }) {
  const params = useParams();
  const grpID = params.id;
  const [errors, setErrors] = useState(null);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [dateTime, setDateTime] = useState("");
  const [error, setError] = useState("");

  async function createEvent(e) {
    e.preventDefault();

    const Data = {
      title: title,
      description: description,
      dateTime: dateTime
    }

    if (!title || !description || !dateTime) {
      setErrors("All fields are required");
      return;
    }

    if ((title.length < 5 && title.length > 50)) {
      setErrors(" Title must be between 5 and 50 characters");
      return;
    }

    if (description.length < 10 || description.length > 300) {
      setErrors("Description must be between 10 and 300 characters");
      return;
    }


    if (new Date(dateTime) <= new Date()) {
      setErrors("Event date and time must be in the future");
      return;
    }

    try {
      const res = await fetch(`http://localhost:8080/api/createEvent/${grpID}`, {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(Data),
      });

      if (!res.ok) {
        let err = await res.text();
        setErrors(err);
        return;
      }

      await fetchEvents();

      closeForm();

      setErrors(null);




    } catch (error) {
      setErrors(error.message);
    }


  }



  return (

    <div className="events-container">
      <div className="create-event-card">
        <h2 className="card-title">Create Event</h2>
        <form className="event-form">
          <div className="form-group">
            <label htmlFor="event-title">Title</label>
            <input type="text" id="event-title" placeholder="Event title..." onChange={(e) => { setTitle(e.target.value) }} />
          </div>

          <div className="form-group">
            <label htmlFor="event-description">Description</label>
            <textarea id="event-description" rows="3" placeholder="Event description..." onChange={(e) => { setDescription(e.target.value) }}></textarea>
          </div>

          <div className="form-group">
            <label htmlFor="event-datetime">Day/Time</label>
            <input type="datetime-local" id="event-datetime" onChange={(e) => { setDateTime(e.target.value) }} />
          </div>
          <span className="error" style={{ red: "red" }}> {errors}</span>
          <button type="submit" className="btn-create" onClick={createEvent}>Create Event</button>
        </form>
        <p className="error-message" style={{ color: "red", alignItems: 'center' }}>{error}</p>
      </div>


    </div>
  )

}
export function GroupChat() {
  const [id, setId] = useState(null)
  const { addEventListener, removeEventListener, sendMessage } = useWS();

  useEffect(() => {
    
    addEventListener("group_message", handleGroupChatMessage);
    return () => {
      removeEventListener("group_message", handleGroupChatMessage);
    };
  }, [addEventListener, removeEventListener]);
  useEffect(() => {
    fetch("http://localhost:8080/api/me")
      .then(res => res.json())
      .then(data => setId(data.user_id))
      .catch(err => console.error(err))
  })
  function handleSendGroupChatMessage(input) {
    if (input.trim() === "") return
    const payload = {
      senderId: id,
      messageContent: input,
      type: "message",
    }
    sendMessage(payload)

  }
  return (
    <div className="group-empty-state">
      <MessageCircle className="tab-icon" />
      <p className="group-empty-state-text">Chat interface coming soon</p>
    </div>
  );
}
