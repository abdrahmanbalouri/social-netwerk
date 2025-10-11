"use client";

import LeftBar from "../../../components/LeftBar";
import RightBar from '../../../components/RightBar';
import Navbar from '../../../components/Navbar';
import { useDarkMode } from '../../../context/darkMod.js';


export default function chat() {
    const { darkMode } = useDarkMode();

    return (
        <div className={darkMode ? 'theme-dark' : 'theme-light'}>
            <Navbar/>
            <main className="content">
            <div><LeftBar /></div>
                    <h1>jmii</h1>
                <div><RightBar /></div>
            </main>

        </div>
    )
}