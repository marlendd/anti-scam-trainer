import {Link} from "react-router-dom";
import {Icon} from "@/shared/ui/icon";
import {faArrowRightFromBracket} from "@fortawesome/free-solid-svg-icons";
import {useLogout} from "@/features/logout-button";

export function LogoutButton() {

    const logout = useLogout();

    const handleLogout = () => {
        logout.mutate();
    }

    return (
        <Link to="/" onClick={handleLogout}>
            <span>Выйти</span>
            <Icon icon={faArrowRightFromBracket} style={{color: 'inherit', height: '14px'}}/>
        </Link>
    )
}