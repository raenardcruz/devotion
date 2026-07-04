import os
import json

def generate_index():
    resources_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), '../src/components/resources'))
    output_file = os.path.join(resources_dir, 'index.json')
    
    entries = []
    
    for filename in os.listdir(resources_dir):
        if filename.endswith('.md'):
            filepath = os.path.join(resources_dir, filename)
            doc_id = os.path.splitext(filename)[0]
            
            # Default title in case we can't find one
            title = doc_id.replace('_', ' ').title()
            
            try:
                with open(filepath, 'r', encoding='utf-8') as f:
                    for line in f:
                        line = line.strip()
                        if line.startswith('# '):
                            title = line[2:].strip()
                            break
            except Exception as e:
                print(f"Error reading {filename}: {e}")
                
            entries.append({
                "id": doc_id,
                "title": title
            })
            
    # Sort entries alphabetically by title
    entries.sort(key=lambda x: x['title'].lower())
    
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(entries, f, indent=2, ensure_ascii=False)
        
    print(f"Generated index with {len(entries)} entries at {output_file}")

if __name__ == '__main__':
    generate_index()
