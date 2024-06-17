import { ReactNode, createContext, useContext, useState } from "react";
import { AccountInfo, DrawInfo, NokInfo, TaskInfo } from "src/services/user.d";

export type UserInfo = {
  accountInfo: AccountInfo;
  drawInfo: DrawInfo[];
  nokInfo: NokInfo;
  taskInfos: TaskInfo[];
  metaInfo: any;
};

export interface UserContextProps {
  userData: UserInfo | null;
  setUserData: (user: UserInfo | null) => void;
}

export const UserContext = createContext<UserContextProps | undefined>(
  undefined
);

export const UserProvider: React.FC<{ children: ReactNode }> = ({
  children,
}) => {
  const [user, setUser] = useState<UserInfo | null>(null);

  return (
    <UserContext.Provider value={{ userData: user, setUserData: setUser }}>
      {children}
    </UserContext.Provider>
  );
};

export const useUserInfo = (): UserContextProps => {
  const context = useContext(UserContext);

  if (context === undefined) {
    throw new Error("must be used within a UserProvider");
  }

  return context;
};
