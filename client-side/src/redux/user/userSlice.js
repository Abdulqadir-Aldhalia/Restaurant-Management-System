import { createSlice } from "@reduxjs/toolkit";
const initialState = {
  userToken: localStorage.getItem("user_token") ?? "",
};
export const userSlice = createSlice({
  name: "user", // name identify the slice
  initialState, // inital state
  reducers: {
    // use this action when user logged in
    setUserToken: (state, action) => {
      state.userToken = action.payload;
      localStorage.setItem("user_token", action.payload.toString());
    },
    // use this function when user logged out or token expired
    removeUserToken: (state) => {
      state.userToken = "";
      localStorage.removeItem("user_token");
    },
  },
});

// Action creators are generated for each case reducer function
export const { setUserToken, removeUserToken } = userSlice.actions;

export default userSlice.reducer;

