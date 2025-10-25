"use client"
import LeftBar from '../../components/LeftBar';
import Navbar from '../../components/Navbar';
import RightBar from '../../components/RightBar';
import { useDarkMode } from '../../context/darkMod';

const page = () => {
    const { darkMode } = useDarkMode();

    return (
        <div id="div" className={darkMode ? 'theme-dark' : 'theme-light'}>
            <Navbar />
            <main className="content">
                <LeftBar showSidebar={true} />
                <div className="game-container">
                    <h1 className="game-title">Watch Page</h1>
                    <p>This is the watch page where videos will be displayed.</p>
                </div>
                <RightBar />
            </main>
        </div>
    )
}

export default page
