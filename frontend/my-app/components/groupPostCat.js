import {useState} from 'react'
import { AllPosts } from "../app/groups/[id]/page";

export function GroupPostChat() {
    const [activeTab, setActiveTab] = useState('posts');
    return (
        <div>
            <div className='tabs-container '>
                <button className={activeTab === 'posts' ? 'group-tab-button-group active' : 'group-tab-button-group'}
                    onClick={() => setActiveTab('posts')}>Posts</button>
                <button className={activeTab === 'chat' ? 'group-tab-button-group active' : 'group-tab-button-group'}
                    onClick={() => setActiveTab('chat')}>Chat</button>
            </div>

            {/* <div className="group-container"> */}
            {activeTab === 'posts' ? (
                <AllPosts />
            ) : (
                console.log("hhhhhhhhh")
            )}
            {/* </div> */}
        </div>
    );
}