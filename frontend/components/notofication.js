import Link from "next/link";
import "../styles/notifaction.css";
import { useDarkMode } from "../context/darkMod";

export default function Notification({ data, onClose }) {

    const { darkMode } = useDarkMode();


    return (
        <div className={`notification ${darkMode ? 'theme-dark' : 'theme-light'}`}>
            <div className="notification-header">
                <div style={{ display: 'flex', alignItems: 'center' }}>
                    <h3 className="notification-title">New notification</h3>
                </div>
                <button
                    className="notification-close"
                    onClick={onClose}
                    aria-label="Close notification"
                >
                    <i className="fa fa-times"></i>
                </button>
            </div>

            <Link href={`/profile/${data.from}`}>
                <div className="notification-container">
                    <div className="notification-media">
                        <img
                            src={data?.photo ? `/uploads/${data.photo}` : "/assets/default.png"}
                            alt={`${data?.name || 'User'} profile picture`}
                            className="notification-user-avatar"
                        />
                    </div>

                    <div className="notification-content">
                        <p className="notification-text">
                            <strong>{data?.name}</strong> {data?.content}
                        </p>

                        <span className="notification-timer">
                            {'a few seconds ago'}
                        </span>
                    </div>
                </div>
            </Link>
        </div>
    );
}
