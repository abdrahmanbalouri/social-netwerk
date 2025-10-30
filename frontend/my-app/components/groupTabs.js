"use client"
import { useState } from 'react'
// import "../styles/groupstyle.css"
import { MyGroups, AllGroups } from "../app/groups/page"

export function GroupsTabs() {
    const [activeTab, setActiveTab] = useState('my-groups');
    return (
        <div className='tabs-container'>
            <div className="tabs">
                <button className={activeTab === 'my-groups' ? 'tab-buttonn active' : 'tab-buttonn'}
                    onClick={() => setActiveTab('my-groups')}>My Groups</button>
                <button className={activeTab === 'all-groups' ? 'tab-buttonn active' : 'tab-buttonn'}
                    onClick={() => setActiveTab('all-groups')}>All Groups</button>
            </div>
            <div className="group-container">
                {activeTab === 'my-groups' ? (
                    <MyGroups />
                ) : (
                    <AllGroups />
                )}
            </div>
        </div>
    );
}