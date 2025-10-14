"use client";
import { fips } from "crypto";
import Navbar from "../../../components/Navbar.js"
// import { useRouter } from 'next/navigation'
import Post from "../../../components/Post.js";
import { useEffect, useState } from "react";


export default function () {
    return (
        <>
            <Navbar />
            <AllPosts />
        </>
    )
}

function AllPosts() {
    const [posts, setPost] = useState(null);
    const [loading, setLoading] = useState(true);
    
    
    useEffect(() => {
        console.log(121212121212);
        const grpID = localStorage.getItem('selectedGroup')
        console.log("Group id is :22222222", grpID);
        if (!grpID) return
        fetch('http://localhost:8080/group/fetchPosts', {
            method: 'POST',
            credentials: 'include',
            body: JSON.stringify({
                grpId: grpID
            }),
            credentials: 'include',
        })
            .then(res => res.json())
            .then(data => {
                console.log("data is :", data);
                setPost(data);
                setLoading(false);
            })
            .catch(error => {
                console.error("Failed to fetch posts for data:", error);
                setLoading(false);
            });
    }, [])
    if (!posts) {
        return <div>There is no post yetttttt.</div>;
    }
    console.log("posts are :", posts);
    return (
        posts.map((post) => (
            <Post
                key={post.id}
                post={post}
                onGetComments={GetComments}
                ondolike={Handlelik}
            />
        ))
    )
}

function GetComments(post) {
    // setSelectedPost({
    //     id: post.id,
    //     title: post.title || post.post_title || "Post"
    // });
}

function Handlelik() {

}