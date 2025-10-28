// import "../styles/groupstyle.css"
export function GroupCard({ group, onShow }) {
    console.log("inside groupCard function/component");
    return (
        <div className="group-card">
            <div className="group-content">
                <h2 className="group-title">{group.Title}</h2>
                <p className="group-description">{group.Description}</p>
            </div>
            <div className="group-footer">
                <button onClick={() => onShow(group.ID)}>Show</button>
            </div>
        </div>
    );
}
