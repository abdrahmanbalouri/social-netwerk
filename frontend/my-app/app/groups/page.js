"use client"
import Navbar from '../../components/Navbar.js';
import { useEffect, useState } from "react";
import { useRouter } from 'next/router';
import "./page.css"


export default function () {
    return (
        <>
            <Navbar />
            {/* <Groups /> */}
            <GroupsTabs />
        </>
    )
}
function GroupsTabs() {
    const [activeTab, setActiveTab] = useState('my-groups');
    return (
        <>
            <div className="tabs-container">
                <button className={activeTab === 'my-groups' ? 'tab-button active' : 'tab-button'}
                    onClick={() => setActiveTab('my-groups')}>My Groups</button>
                <button className={activeTab === 'all-groups' ? 'tab-button active' : 'tab-button'}
                    onClick={() => setActiveTab('all-groups')}>All Groups</button>
            </div>
        </>
    );
}



