'use client';
import { useState, useEffect, useRef } from 'react';
import Link from 'next/link';
import "../../styles/watch.css";

export default function Reels() {
    const [videos, setVideos] = useState([]);
    const [currentIndex, setCurrentIndex] = useState(0);
    const [loading, setLoading] = useState(true);

    // const handleGoHome = () => {
    //     window.location.href = "/home";
    // };

    useEffect(() => {
        const fetchVideos = async () => {
            try {
                const res = await fetch('http://localhost:8080/api/videos', {
                    method: 'GET',
                    credentials: 'include',
                    headers: { 'Content-Type': 'application/json' },
                });
                if (!res.ok) return;
                const data = await res.json();
                setVideos(data);
            } catch (error) {
                console.error("Failed to fetch videos:", error);
            } finally {
                setLoading(false);
            }
        };
        fetchVideos();
    }, []);

 
    if (loading) return <div className="loading">Loading...</div>;
    if (!videos || videos.length === 0) return (
        <div className="no-videos">
            No reels
            <Link
                href="/home"
                className="home-link"
                style={{
                    marginTop: '20px',
                    display: 'inline-block',
                    padding: '10px 20px',
                    color: '#fff',
                    backgroundColor: '#5271ff',
                    borderRadius: '5px',
                    textDecoration: 'none',
                }}
            >
                Back to Home
            </Link>
        </div>
    );

    const current = videos[currentIndex];

    const goNext = () =>  setCurrentIndex(i => i + 1);
    const goPrev = () => setCurrentIndex(i => i - 1);

    return (
        <div
            className="reels-container"
          
        >
            {/* Video */}
            <div className="video-wrapper">
                <video
                    key={current.id}
                    src={current.image_path}
                    muted
                    loop
                    playsInline
                    autoPlay
                    className="reel-video"
                />
            </div>

            <div className="like-count">
                <span className="icon">❤️</span>
                <span className="count">{current.like || 0}</span>
            </div>

            {/* Next / Prev Buttons */}
            <div className="controls">
                {currentIndex > 0 && (
                    <button onClick={goPrev} className="nav-btn">Previous</button>
                )}
                <span className="counter">{currentIndex + 1} / {videos.length}</span>
                {currentIndex < videos.length - 1 && (
                    <button onClick={goNext} className="nav-btn">Next</button>
                )}
            </div>

            {/* Progress Bar */}
            <div className="progress-bar">
                <div
                    className="progress-fill"
                    style={{ width: `${((currentIndex + 1) / videos.length) * 100}%` }}
                />
            </div>

            {/* Back to Home */}
            <Link href='/home' className='home-btn' > Go Home</Link>

        </div>
    );
}