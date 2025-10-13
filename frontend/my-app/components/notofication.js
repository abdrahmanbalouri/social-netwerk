import Link from "next/link";
import "../styles/notifaction.css";

export default function Notification({ data }) {
    if (!data) return null;
    console.log();


    return (
        <div className="notification">
            <div className="notification-header">
                <h3 className="notification-title">New notification</h3>
                <i className="fa fa-times notification-close"></i>
            </div>
            <div className="notification-container">
                <div className="notification-media">
                    <img src={data?.photo ? `/uploads/${data.photo}` : "/uploads/default.png"} alt="" className="notification-user-avatar" />
                    <i className="fa fa-thumbs-up notification-reaction"></i>
                </div>
                <div className="notification-content">
                    <p className="notification-text">
                        <strong>{data && data.name}</strong>,  {data && data.messageContent} 
                    </p>
                    <span className="notification-timer">a few seconds ago</span>
                </div>
                <span className="notification-status"></span>
            </div>
        </div>
    );
}
