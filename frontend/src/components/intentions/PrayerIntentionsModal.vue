<template>
  <transition name="fade">
    <div 
      v-if="isOpen" 
      class="fixed inset-0 bg-[#322D29]/50 backdrop-blur-sm z-[100] flex items-center justify-center p-3 sm:p-4 overflow-y-auto"
      @click.self="$emit('close')"
    >
      <div 
        class="bg-gradient-to-b from-white via-[#FAF7F2] to-[#EFE9E1] border border-[#D1C7BD] w-full max-w-xl rounded-[2rem] shadow-2xl flex flex-col max-h-[90vh] overflow-hidden transition-all duration-300 transform scale-100 my-auto"
      >
        <!-- Modal Header -->
        <div class="p-5 border-b border-[#D1C7BD]/60 flex items-center justify-between bg-white/70 backdrop-blur-md">
          <div class="flex items-center space-x-3">
            <div class="w-10 h-10 rounded-2xl bg-gradient-to-br from-[#72383D]/15 to-[#AC9C8D]/20 border border-[#D1C7BD]/80 flex items-center justify-center text-[#72383D] text-lg shadow-xs">
              🙏
            </div>
            <div>
              <div class="flex items-center space-x-2">
                <h3 class="font-serif text-lg md:text-xl text-[#72383D] font-bold">
                  Prayer Intentions
                </h3>
                <span class="px-2.5 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-[#72383D]/10 text-[#72383D] border border-[#72383D]/20">
                  {{ activeCount }} Active
                </span>
              </div>
              <p class="text-[11px] text-[#322D29]/60 mt-0.5">
                {{ devotionTitle }} &bull; Offered after the Sign of the Cross
              </p>
            </div>
          </div>

          <button 
            @click="$emit('close')" 
            class="w-8 h-8 rounded-full border border-[#D1C7BD] flex items-center justify-center text-[#322D29]/60 hover:bg-[#D1C7BD]/30 hover:text-[#72383D] active:scale-95 transition-all p-0 cursor-pointer"
            aria-label="Close"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <!-- Scrollable Content -->
        <div class="p-4 sm:p-5 overflow-y-auto space-y-4 flex-grow">
          
          <!-- Todo Add Form Card -->
          <div class="bg-white/80 border border-[#D1C7BD]/80 rounded-2xl p-3.5 sm:p-4 shadow-xs space-y-3">
            <div class="flex items-center justify-between">
              <span class="text-[11px] font-bold text-[#72383D] uppercase tracking-wider flex items-center gap-1.5">
                <span>➕</span>
                <span>Add Intention</span>
              </span>
              <button 
                type="button"
                @click="showPresets = !showPresets"
                class="text-[10px] font-bold text-[#72383D] hover:underline flex items-center gap-1 cursor-pointer"
              >
                <span>{{ showPresets ? 'Hide Suggestions' : '💡 Suggestions' }}</span>
              </button>
            </div>

            <!-- Quick Template Chips -->
            <div v-if="showPresets" class="flex flex-wrap gap-1.5 pt-1 pb-2 animate-fade-in-down">
              <button
                v-for="(preset, pIdx) in INTENTION_PRESETS"
                :key="pIdx"
                @click="applyPreset(preset)"
                class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[11px] bg-white border border-[#D1C7BD]/80 hover:border-[#72383D]/60 hover:bg-[#72383D]/5 text-[#322D29] transition-all cursor-pointer shadow-2xs"
              >
                <span>{{ preset.icon }}</span>
                <span class="truncate max-w-[200px]">{{ preset.title }}</span>
              </button>
            </div>

            <!-- Input Fields -->
            <div class="space-y-2.5">
              <!-- Title Input & Category Dropdown Trigger -->
              <div class="flex flex-col sm:flex-row gap-2">
                <div class="relative flex-grow">
                  <input 
                    type="text" 
                    v-model="newTitle" 
                    @keydown.enter.prevent="handleAddIntention"
                    placeholder="E.g., For my mother's healing..."
                    class="w-full px-3.5 py-2 bg-white border border-[#D1C7BD] rounded-xl text-xs sm:text-sm text-[#322D29] placeholder-[#322D29]/40 focus:outline-none focus:ring-1 focus:ring-[#72383D]/40 focus:border-[#72383D] transition-all"
                  />
                </div>

                <!-- Category / Icon Selector -->
                <div class="relative">
                  <button 
                    type="button"
                    @click="showCategoryDropdown = !showCategoryDropdown"
                    class="w-full sm:w-auto inline-flex items-center justify-between gap-2 px-3 py-2 bg-white border border-[#D1C7BD] rounded-xl text-xs font-semibold text-[#322D29] hover:bg-[#FAF7F2] transition-colors cursor-pointer"
                  >
                    <span class="flex items-center gap-1.5">
                      <span class="text-base">{{ selectedCategory.icon }}</span>
                      <span class="truncate max-w-[110px]">{{ selectedCategory.name }}</span>
                    </span>
                    <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="text-[#322D29]/50">
                      <path d="m6 9 6 6 6-6"/>
                    </svg>
                  </button>

                  <!-- Category Dropdown Menu -->
                  <div 
                    v-if="showCategoryDropdown" 
                    class="absolute right-0 top-full mt-1.5 z-50 w-64 bg-white/95 backdrop-blur-md border border-[#D1C7BD] rounded-xl shadow-xl p-1.5 space-y-1 max-h-56 overflow-y-auto"
                  >
                    <button 
                      v-for="cat in PRAYER_CATEGORIES" 
                      :key="cat.id"
                      type="button"
                      @click="selectCategory(cat)"
                      class="w-full text-left p-2 rounded-lg text-xs flex items-center gap-2 hover:bg-[#72383D]/10 transition-colors cursor-pointer"
                      :class="[selectedCategory.id === cat.id ? 'bg-[#72383D]/10 font-bold text-[#72383D]' : 'text-[#322D29]']"
                    >
                      <span class="text-base">{{ cat.icon }}</span>
                      <div class="flex-grow min-w-0">
                        <div class="font-medium truncate">{{ cat.name }}</div>
                        <div class="text-[9px] text-[#322D29]/50 truncate">{{ cat.description }}</div>
                      </div>
                    </button>
                  </div>
                </div>
              </div>

              <!-- Description / Notes (Optional) -->
              <div>
                <textarea 
                  v-model="newDescription"
                  rows="2"
                  placeholder="Description or prayer notes (optional)..."
                  class="w-full px-3.5 py-2 bg-white border border-[#D1C7BD] rounded-xl text-xs text-[#322D29] placeholder-[#322D29]/40 focus:outline-none focus:ring-1 focus:ring-[#72383D]/40 focus:border-[#72383D] transition-all resize-none"
                ></textarea>
              </div>

              <!-- Submit Button -->
              <div class="flex items-center justify-end">
                <button 
                  type="button"
                  @click="handleAddIntention"
                  :disabled="!newTitle.trim()"
                  class="px-4 py-2 bg-gradient-to-r from-[#72383D] to-[#322D29] text-white rounded-xl text-xs font-bold uppercase tracking-wider hover:from-[#8B464C] hover:to-[#453E38] active:scale-95 disabled:opacity-40 disabled:cursor-not-allowed transition-all shadow-xs flex items-center gap-1.5 cursor-pointer"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                    <line x1="12" y1="5" x2="12" y2="19"></line>
                    <line x1="5" y1="12" x2="19" y2="12"></line>
                  </svg>
                  <span>Add to List</span>
                </button>
              </div>
            </div>
          </div>

          <!-- Todo List Section Header & Filters -->
          <div class="flex items-center justify-between pt-1">
            <div class="flex items-center space-x-2">
              <span class="text-xs font-serif font-bold text-[#322D29]">Your Intentions List</span>
              <span class="text-[10px] text-[#322D29]/50 font-semibold">({{ intentions.length }})</span>
            </div>
            
            <div class="flex items-center space-x-1.5">
              <button 
                v-if="hasCompleted"
                @click="clearCompleted"
                class="text-[10px] text-[#72383D] hover:underline font-semibold cursor-pointer"
              >
                Clear completed
              </button>
            </div>
          </div>

          <!-- Empty State -->
          <div 
            v-if="intentions.length === 0" 
            class="p-8 rounded-2xl border border-dashed border-[#D1C7BD] bg-white/40 text-center flex flex-col items-center justify-center space-y-2"
          >
            <div class="w-12 h-12 rounded-full bg-[#EFE9E1] flex items-center justify-center text-xl shadow-2xs">
              🕊️
            </div>
            <div class="font-serif text-sm font-bold text-[#72383D]">No intentions added yet</div>
            <p class="text-xs text-[#322D29]/60 max-w-xs leading-relaxed">
              Add your personal intentions above. They will be presented right after the opening Sign of the Cross during your prayer.
            </p>
          </div>

          <!-- Intentions List (Todo Style) -->
          <div v-else class="space-y-2">
            <div 
              v-for="item in intentions" 
              :key="item.id"
              class="border rounded-2xl transition-all duration-200 shadow-2xs overflow-hidden"
              :class="[
                item.completed 
                  ? 'bg-white/40 border-[#D1C7BD]/50 opacity-60' 
                  : 'bg-white/85 border-[#D1C7BD] hover:border-[#72383D]/40'
              ]"
            >
              <!-- Editing Mode Form -->
              <div v-if="editingId === item.id" class="p-3.5 bg-white space-y-3">
                <div class="flex items-center justify-between border-b border-[#D1C7BD]/50 pb-2">
                  <span class="text-[10px] font-bold uppercase tracking-wider text-[#72383D]">Edit Intention</span>
                  <button @click="cancelEdit" class="text-xs text-[#322D29]/50 hover:text-[#72383D] cursor-pointer">Cancel</button>
                </div>

                <div class="space-y-2">
                  <div class="flex flex-col sm:flex-row gap-2">
                    <input 
                      type="text" 
                      v-model="editTitle" 
                      class="flex-grow px-3 py-1.5 bg-white border border-[#D1C7BD] rounded-lg text-xs text-[#322D29] focus:outline-none focus:ring-1 focus:ring-[#72383D]"
                    />
                    
                    <select 
                      v-model="editCategory" 
                      class="px-2.5 py-1.5 bg-white border border-[#D1C7BD] rounded-lg text-xs text-[#322D29] focus:outline-none"
                    >
                      <option v-for="cat in PRAYER_CATEGORIES" :key="cat.id" :value="cat.id">
                        {{ cat.icon }} {{ cat.name }}
                      </option>
                    </select>
                  </div>

                  <textarea 
                    v-model="editDescription" 
                    rows="2"
                    placeholder="Description / notes..."
                    class="w-full px-3 py-1.5 bg-white border border-[#D1C7BD] rounded-lg text-xs text-[#322D29] focus:outline-none focus:ring-1 focus:ring-[#72383D] resize-none"
                  ></textarea>

                  <div class="flex justify-end gap-2 pt-1">
                    <button 
                      @click="cancelEdit" 
                      class="px-3 py-1 border border-[#D1C7BD] text-xs font-semibold rounded-lg hover:bg-black/5 cursor-pointer"
                    >
                      Cancel
                    </button>
                    <button 
                      @click="saveEdit(item.id)" 
                      :disabled="!editTitle.trim()"
                      class="px-3.5 py-1 bg-[#72383D] text-white text-xs font-bold rounded-lg hover:bg-[#8B464C] cursor-pointer disabled:opacity-40"
                    >
                      Save Changes
                    </button>
                  </div>
                </div>
              </div>

              <!-- Normal Display Item -->
              <div v-else class="p-3 sm:p-3.5 flex items-start space-x-3">
                <!-- Checkbox (Todo Style Toggle) -->
                <button 
                  type="button"
                  @click="toggleComplete(item.id)" 
                  class="mt-0.5 w-5 h-5 rounded-full border flex items-center justify-center transition-all cursor-pointer shrink-0"
                  :class="[
                    item.completed 
                      ? 'bg-[#30D158] border-[#30D158] text-white' 
                      : 'border-[#D1C7BD] hover:border-[#72383D] bg-white'
                  ]"
                  :title="item.completed ? 'Mark active' : 'Mark completed/offered'"
                >
                  <svg v-if="item.completed" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="20 6 9 17 4 12"></polyline>
                  </svg>
                </button>

                <!-- Icon Badge -->
                <div class="w-8 h-8 rounded-xl bg-[#FAF7F2] border border-[#D1C7BD]/70 flex items-center justify-center text-sm shadow-2xs shrink-0">
                  {{ item.icon || '🙏' }}
                </div>

                <!-- Text Info -->
                <div class="flex-grow min-w-0">
                  <div class="flex items-center gap-2">
                    <h4 
                      class="text-xs sm:text-sm font-bold text-[#322D29] leading-snug"
                      :class="{ 'line-through text-[#322D29]/40': item.completed }"
                    >
                      {{ item.title }}
                    </h4>
                    <span 
                      v-if="getCategoryName(item.category)" 
                      class="text-[9px] font-semibold text-[#72383D] px-1.5 py-0.5 bg-[#72383D]/10 rounded-md shrink-0"
                    >
                      {{ getCategoryName(item.category) }}
                    </span>
                  </div>

                  <!-- Description -->
                  <p 
                    v-if="item.description" 
                    class="text-[11px] text-[#322D29]/70 mt-1 leading-relaxed whitespace-pre-line"
                    :class="{ 'line-through text-[#322D29]/30': item.completed }"
                  >
                    {{ item.description }}
                  </p>
                </div>

                <!-- Item Actions -->
                <div class="flex items-center space-x-1 shrink-0 ml-1">
                  <button 
                    @click="startEdit(item)" 
                    class="p-1.5 rounded-lg text-[#322D29]/50 hover:text-[#72383D] hover:bg-[#72383D]/10 active:scale-95 transition-all cursor-pointer"
                    title="Modify intention"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"></path>
                      <path d="m15 5 4 4"></path>
                    </svg>
                  </button>

                  <button 
                    @click="deleteIntention(item.id)" 
                    class="p-1.5 rounded-lg text-[#322D29]/40 hover:text-red-600 hover:bg-red-50 active:scale-95 transition-all cursor-pointer"
                    title="Delete intention"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M3 6h18"></path>
                      <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"></path>
                      <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"></path>
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Helpful note -->
          <div class="bg-[#EFE9E1]/80 border border-[#D1C7BD] rounded-xl p-3 flex items-start space-x-2.5 text-[#322D29]/70">
            <svg class="text-[#72383D] mt-0.5 shrink-0" xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="16" x2="12" y2="12"></line>
              <line x1="12" y1="8" x2="12.01" y2="8"></line>
            </svg>
            <p class="text-[10px] leading-relaxed">
              When praying, your active intentions will be displayed directly after the opening Sign of the Cross so you can lift them up to Our Lord.
            </p>
          </div>

        </div>

        <!-- Modal Footer -->
        <div class="p-4 bg-white/90 border-t border-[#D1C7BD]/60 flex items-center justify-between">
          <div class="text-[11px] text-[#322D29]/60">
            <span class="font-bold text-[#72383D]">{{ activeCount }}</span> of {{ intentions.length }} active
          </div>
          
          <button 
            @click="$emit('close')" 
            class="px-6 py-2 bg-gradient-to-r from-[#72383D] to-[#322D29] hover:from-[#8B464C] hover:to-[#453E38] text-white rounded-full text-xs font-bold uppercase tracking-wider outline-none shadow-xs transition-all cursor-pointer"
          >
            Done
          </button>
        </div>

      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { 
  usePrayerIntentions, 
  PRAYER_CATEGORIES, 
  INTENTION_PRESETS,
  type DevotionKey, 
  type PrayerCategory, 
  type PrayerIntention 
} from '../../composables/usePrayerIntentions';

const props = defineProps<{
  isOpen: boolean;
  devotionKey: DevotionKey;
  devotionTitle: string;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'updated'): void;
}>();

const { 
  intentions, 
  activeCount, 
  addIntention, 
  updateIntention, 
  deleteIntention: deleteFromStore, 
  toggleComplete: toggleFromStore,
  clearCompleted: clearCompletedFromStore
} = usePrayerIntentions(props.devotionKey);

// New intention form state
const newTitle = ref('');
const newDescription = ref('');
const selectedCategory = ref<PrayerCategory>(PRAYER_CATEGORIES[0]!);
const showCategoryDropdown = ref(false);
const showPresets = ref(false);

// Editing state
const editingId = ref<string | null>(null);
const editTitle = ref('');
const editDescription = ref('');
const editCategory = ref('');

const hasCompleted = computed(() => intentions.value.some(i => i.completed));

function selectCategory(cat: PrayerCategory) {
  selectedCategory.value = cat;
  showCategoryDropdown.value = false;
}

function applyPreset(preset: typeof INTENTION_PRESETS[0]) {
  newTitle.value = preset.title;
  newDescription.value = preset.description;
  const foundCat = PRAYER_CATEGORIES.find(c => c.id === preset.category);
  if (foundCat) {
    selectedCategory.value = foundCat;
  }
  showPresets.value = false;
}

function handleAddIntention() {
  if (!newTitle.value.trim()) return;

  addIntention({
    title: newTitle.value,
    description: newDescription.value,
    category: selectedCategory.value.id,
    icon: selectedCategory.value.icon
  });

  newTitle.value = '';
  newDescription.value = '';
  showCategoryDropdown.value = false;
  emit('updated');
}

function startEdit(item: PrayerIntention) {
  editingId.value = item.id;
  editTitle.value = item.title;
  editDescription.value = item.description || '';
  editCategory.value = item.category || PRAYER_CATEGORIES[0]!.id;
}

function cancelEdit() {
  editingId.value = null;
  editTitle.value = '';
  editDescription.value = '';
}

function saveEdit(id: string) {
  if (!editTitle.value.trim()) return;

  const foundCat = PRAYER_CATEGORIES.find(c => c.id === editCategory.value) || PRAYER_CATEGORIES[0]!;
  updateIntention(id, {
    title: editTitle.value,
    description: editDescription.value,
    category: foundCat.id,
    icon: foundCat.icon
  });

  cancelEdit();
  emit('updated');
}

function deleteIntention(id: string) {
  deleteFromStore(id);
  emit('updated');
}

function toggleComplete(id: string) {
  toggleFromStore(id);
  emit('updated');
}

function clearCompleted() {
  clearCompletedFromStore();
  emit('updated');
}

function getCategoryName(categoryId: string) {
  const found = PRAYER_CATEGORIES.find(c => c.id === categoryId);
  return found ? found.name : '';
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.animate-fade-in-down {
  animation: fadeInDown 0.3s ease-out forwards;
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
