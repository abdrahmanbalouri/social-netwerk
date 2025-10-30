// import "../styles/groupstyle.css"
import {useState} from 'react'
// import "../styles/grouppage.css"
import { AllPosts } from "../app/groups/[id]/page";

export function GroupPostChat() {
    const [activeTab, setActiveTab] = useState('posts');
    return (
        <div className="main-container">
            <div className='tabs'>
                <button className={activeTab === 'posts' ? 'tab-buttonn active' : 'tab-buttonn'}
                    onClick={() => setActiveTab('posts')}>Posts</button>
                <button className={activeTab === 'chat' ? 'tab-buttonn active' : 'tab-buttonn'}
                    onClick={() => setActiveTab('chat')}>Chat</button>
            </div>

            {activeTab === 'posts' ? (
                <AllPosts />
            ) : (
                <dibv>hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh</dibv>
            )}
            {/* </div> */}
        </div>
    );
}