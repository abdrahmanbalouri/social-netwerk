"use client";
import Navbar from "../../../components/Navbar.js"
import { GroupPostChat } from "../../../components/groupPostCat.js";
import Post from "../../../components/Post.js";
import { useEffect, useState } from "react";
import "../../../styles/groupstyle.css"
import "../../../styles/group.css"
import { useParams } from "next/navigation";
import { PostCreationTrigger } from "../../../components/cretaePostGroup.js"
import LeftBar from "../../../components/LeftBar.js";
import RightBarGroup from "../../../components/RightBarGroups.js";
import { useDarkMode } from "../../../context/darkMod.js";

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
    )
}

function sendRequest(invitedUserID, grpID) {
    console.log("invited user id is : ", invitedUserID);
    console.log("group id is : ", grpID);

    fetch(`http://localhost:8080/group/invitation/${grpID}`, {
        method: 'POST',
        credentials: 'include',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            invitedUsers: [invitedUserID],
        }),
    })
        .then(async res => {
            const temp = await res.json()
            console.log("temp is :", temp);
        })
        .catch(error => {
            console.log("error sending invitation to the user :", error);
        }
        )
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
            method: 'GET',
            credentials: 'include',
        })
            .then(res => {
                if (!res.ok) throw new Error('Failed to fetch posts');
                return res.json();
            })
            .then(data => {
                setPost(data || []);
                setLoading(false);
            })
            .catch(error => {
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

            {posts.length === 0 ? (
                <div>There is no post yeeeeeet.</div>
            ) : (
                <div className="posts-list">
                    {posts.map((post) => (
                        <Post
                            key={post.id}
                            post={post}
                            onGetComments={GetComments}
                            ondolike={AddLike}
                        />
                    ))}
                </div>
            )}
        </div>
    )
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
            <div className="posts-list">
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

function AddLike() {

}

export async function CreatePost(groupId, formData) {

    const res = await fetch(`http://localhost:8080/group/addPost/${groupId}`, {
        method: "POST",
        credentials: "include",
        body: formData,
    });

    const text = await res.text();

    if (!res.ok) throw new Error("Failed to create a post for groups");

    return JSON.parse(text);
}
