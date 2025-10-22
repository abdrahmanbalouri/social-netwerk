import "../styles/groupTabs.css"
import {useState} from 'react'
import {MyGroups, AllGroups} from "../app/groups/page"

export function GroupsTabs() {
    const [activeTab, setActiveTab] = useState('my-groups');
    return (
        <>
            <div className='tabs-container '>
                <button className={activeTab === 'my-groups' ? 'tab-button-group active' : 'tab-button-group'}
                    onClick={() => setActiveTab('my-groups')}>My Groups</button>
                <button className={activeTab === 'all-groups' ? 'tab-button-group active' : 'tab-button-group'}
                    onClick={() => setActiveTab('all-groups')}>All Groups</button>
            </div>

            {/* <div className="group-container"> */}
            {activeTab === 'my-groups' ? (
                <MyGroups />
            ) : (
                <AllGroups />
            )}
            {/* </div> */}
        </>
    );
}