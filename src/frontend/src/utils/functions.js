export function getTimestamp(timestamp) {
    const date = new Date(timestamp);
    const hour = date.getHours();
    const minutes = date.getMinutes();
    const day = date.getDate();
    const monthIndex = date.getMonth();
    const year = date.getFullYear();
    return (hour + ':' + minutes + ' ' + day + '/' + (monthIndex + 1) + ' ' + year);
}