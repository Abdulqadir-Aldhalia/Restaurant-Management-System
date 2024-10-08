import { createSlice } from "@reduxjs/toolkit";
const initialState = {
  userToken: localStorage.getItem("user_token") ?? "",
};
export const userSlice = createSlice({
  name: "user",
  initialState,
  reducers: {
    setUserToken: (state, action) => {
      state.userToken = action.payload;
      localStorage.setItem("user_token", action.payload.toString());
    },
    removeUserToken: (state) => {
      state.userToken = "";
      localStorage.removeItem("user_token");
    },
  },
});

export const { setUserToken, removeUserToken } = userSlice.actions;

export default userSlice.reducer;
