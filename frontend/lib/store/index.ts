import { configureStore } from '@reduxjs/toolkit';
import repositoryVisibilityReducer from './slices/repositoryVisibilitySlice';

export const makeStore = () => {
  return configureStore({
    reducer: {
      repositoryVisibility: repositoryVisibilityReducer,
    },
  });
};

export type AppStore = ReturnType<typeof makeStore>;
export type RootState = ReturnType<AppStore['getState']>;
export type AppDispatch = AppStore['dispatch'];
