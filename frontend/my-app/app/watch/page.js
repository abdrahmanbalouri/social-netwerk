'use client';
import { useState, useEffect, useRef } from 'react';
import Link from 'next/link';
import "../../styles/watch.css";

export default function Reels() {
    const [videos, setVideos] = useState([]);
    const [currentIndex, setCurrentIndex] = useState(0);
    const [loading, setLoading] = useState(true);
    const touchStartY = useRef(0);

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

    // Swipe
    const handleTouchStart = (e) => {
        touchStartY.current = e.touches[0].clientY;
    };

    const handleTouchEnd = (e) => {
        if (!touchStartY.current) return;
        const touchEndY = e.changedTouches[0].clientY;
        const diff = touchStartY.current - touchEndY;

        if (Math.abs(diff) > 50) {
            if (diff > 0 && currentIndex < videos.length - 1) {
                setCurrentIndex(i => i + 1);
            } else if (diff < 0 && currentIndex > 0) {
                setCurrentIndex(i => i - 1);
            }
        }
        touchStartY.current = 0;
    };

    if (loading) return <div className="loading">Loading...</div>;
    if (videos.length === 0) return (
        <div className="no-videos">
            No reels
            <Link href="/" className="home-link">Back to Home</Link>
        </div>
    );

    const current = videos[currentIndex];

    const goNext = () => currentIndex < videos.length - 1 && setCurrentIndex(i => i + 1);
    const goPrev = () => currentIndex > 0 && setCurrentIndex(i => i - 1);

    return (
        <div
            className="reels-container"
            onTouchStart={handleTouchStart}
            onTouchEnd={handleTouchEnd}
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

            {/*Name + Title*/}
            {/* <div className="user-info">
                <div className="user-text">
                    <p className="name">{current.first_name} {current.last_name}</p>
                    <p className="title">{current.title || ''}</p>
                </div>
            </div> */}

            {/* Like Count */}
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
        <Link href= '/' className='home-btn' > Go Home</Link>
            
        </div>
    );
}