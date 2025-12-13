import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface RepositoryVisibilityState {
  selectedRepoIds: string[];
  isInitialized: boolean;
}

const STORAGE_KEY = 'devplus_selected_repos';

// Load initial state from localStorage
const loadFromStorage = (): string[] => {
  if (typeof window === 'undefined') return [];
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    return stored ? JSON.parse(stored) : [];
  } catch {
    return [];
  }
};

// Save to localStorage
const saveToStorage = (repoIds: string[]) => {
  if (typeof window === 'undefined') return;
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(repoIds));
  } catch (error) {
    console.error('Failed to save repository visibility to localStorage:', error);
  }
};

const initialState: RepositoryVisibilityState = {
  selectedRepoIds: loadFromStorage(),
  isInitialized: false,
};

const repositoryVisibilitySlice = createSlice({
  name: 'repositoryVisibility',
  initialState,
  reducers: {
    initializeSelectedRepos: (state, action: PayloadAction<string[]>) => {
      // Only initialize if not already done and no stored data exists
      if (!state.isInitialized && state.selectedRepoIds.length === 0) {
        state.selectedRepoIds = action.payload;
        state.isInitialized = true;
        saveToStorage(action.payload);
      } else if (!state.isInitialized) {
        state.isInitialized = true;
      }
    },
    toggleRepoVisibility: (state, action: PayloadAction<string>) => {
      const repoId = action.payload;
      const index = state.selectedRepoIds.indexOf(repoId);
      
      if (index > -1) {
        state.selectedRepoIds.splice(index, 1);
      } else {
        state.selectedRepoIds.push(repoId);
      }
      
      saveToStorage(state.selectedRepoIds);
    },
    setAllReposVisibility: (state, action: PayloadAction<string[]>) => {
      state.selectedRepoIds = action.payload;
      saveToStorage(action.payload);
    },
    selectAllRepos: (state, action: PayloadAction<string[]>) => {
      state.selectedRepoIds = action.payload;
      saveToStorage(action.payload);
    },
    deselectAllRepos: (state) => {
      state.selectedRepoIds = [];
      saveToStorage([]);
    },
  },
});

export const {
  initializeSelectedRepos,
  toggleRepoVisibility,
  setAllReposVisibility,
  selectAllRepos,
  deselectAllRepos,
} = repositoryVisibilitySlice.actions;

export default repositoryVisibilitySlice.reducer;
