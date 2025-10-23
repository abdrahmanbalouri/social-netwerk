"use client";
import React, { useEffect } from "react";
import Navbar from "../../components/Navbar.js";
import LeftBar from "../../components/LeftBar.js";
import RightBar from "../../components/RightBar.js";
import { useDarkMode } from "../../context/darkMod.js";
import Link from "next/link.js";
import "../../styles/games.css";
import { useWS } from "../../context/wsContext.js";
import { middleware } from "../../middleware/middelware.js";
import { useRouter } from "next/navigation";

export default function Game() {
  const router = useRouter();
  const sendMessage = useWS()
  // Authentication check
  useEffect(() => {
    const checkAuth = async () => {
      const auth = await middleware();
      if (!auth) {
        router.push("/login");
        sendMessage({ type: "logout" })
      }
    }
    checkAuth();
  }, [])
  const { darkMode } = useDarkMode();
  const [showSidebar, setShowSidebar] = React.useState(true);
  return (
    <div className={darkMode ? "theme-dark" : "theme-light"}>
      <Navbar
        onToggleSidebar={() => setShowSidebar(!showSidebar)}
      />
      <main className="content">
        <LeftBar showSidebar={showSidebar} />

        <div className="game-container">
          <h2 className="game-title">Mini Games</h2>
          <p className="game-desc">Enjoy a selection of simple games below!</p>
          <div className="games-list">
            <div className="game-card">
              <h3>Boomber man</h3>
              <Link href="https://raji383.github.io/bomberman/" target="_blank">
                <img src="/assets/bomberman.png" alt="Boomber man" />
              </Link>
            </div>
            <div className="game-card">
              <h3>Arkanoid</h3>
              <Link href="https://brick-breaker-abaid.netlify.app/" target="_blank">
                <img src="/assets/Arkanoid.png" alt="Arkanoid" />
              </Link>
            </div>
            <div className="game-card">
              <h3>super-mario</h3>
              <Link href="https://raji383.github.io/super-mario/" target="_blank">
                <img src="/assets/super-mario.png" alt="Arkanoid" />
              </Link>
            </div>
            <div className="game-card">
              <h3>Raven</h3>
              <Link href="https://raji383.github.io/Raven-game/" target="_blank">
                <img src="/assets/Raven.png" alt="Arkanoid" />
              </Link>
            </div>
            <div className="game-card">
              <h3>shadow</h3>
              <Link href="https://raji383.github.io/shadow-game/" target="_blank">
                <img src="/assets/shadow.png" alt="Arkanoid" />
              </Link>
            </div>
            <div className="game-card">
              <h3>War-machine</h3>
              <Link href="https://raji383.github.io/War-machine/" target="_blank">
                <img src="/assets/war.png" alt="Arkanoid" />
              </Link>
            </div>
            <div className="game-card">
              <h3>toad</h3>
              <Link href="https://raji383.github.io/egg-game/" target="_blank">
                <img src="/assets/toad.png" alt="Arkanoid" />
              </Link>
            </div>
          </div>
        </div>

        <RightBar />
      </main>
    </div>
  );
}


