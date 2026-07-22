import os
import re
import json
import urllib.request
from html.parser import HTMLParser

class FathersParser(HTMLParser):
    def __init__(self):
        super().__init__()
        self.fathers = []
        self.in_p = False
        self.in_strong = False
        self.in_work_link = False
        self.father_name_parts = []
        self.current_work_href = None
        self.work_title_parts = []
        self.works = []
        self.father_href = None

    def handle_starttag(self, tag, attrs):
        attrs_dict = dict(attrs)
        if tag == 'p':
            self.in_p = True
            self.father_name_parts = []
            self.works = []
            self.father_href = None
        elif tag == 'strong' and self.in_p:
            self.in_strong = True
        elif tag == 'a' and self.in_p:
            href = attrs_dict.get('href', '')
            if '../cathen/' in href:
                self.father_href = href
            elif '../fathers/' in href:
                self.current_work_href = href
                self.in_work_link = True
                self.work_title_parts = []

    def handle_endtag(self, tag):
        if tag == 'p':
            self.in_p = False
            father_name = "".join(self.father_name_parts).strip()
            # Clean up father name
            father_name = re.sub(r'\s*\[.*?\]', '', father_name)
            father_name = re.sub(r'\s*\(.*?\)', '', father_name)
            father_name = father_name.strip()
            
            if father_name and self.works:
                # Check if we already have this father (avoid duplicates)
                existing = next((f for f in self.fathers if f["name"] == father_name), None)
                if existing:
                    existing["works"].extend(self.works)
                else:
                    self.fathers.append({
                        "name": father_name,
                        "works": self.works,
                        "status": "pending"
                    })
            self.works = []
        elif tag == 'strong':
            self.in_strong = False
        elif tag == 'a':
            if self.in_work_link:
                self.in_work_link = False
                work_title = "".join(self.work_title_parts).strip()
                work_title = re.sub(r'\s+', ' ', work_title)
                # Clean up url relative paths
                url_path = self.current_work_href.replace('..', '')
                full_url = f"https://www.newadvent.org{url_path}"
                self.works.append({
                    "title": work_title,
                    "url": full_url,
                    "status": "pending"
                })

    def handle_data(self, data):
        if self.in_p and self.in_strong:
            self.father_name_parts.append(data)
        elif self.in_p and self.in_work_link:
            self.work_title_parts.append(data)

def main():
    url = "https://www.newadvent.org/fathers/index.html"
    print(f"Fetching index from {url}...")
    try:
        req = urllib.request.Request(url, headers={'User-Agent': 'Mozilla/5.0'})
        with urllib.request.urlopen(req) as response:
            html = response.read().decode('utf-8', errors='ignore')
    except Exception as e:
        print(f"Failed to fetch live page: {e}. Trying local fallback...")
        local_path = "/Users/raenard/.gemini/antigravity/brain/4829f76b-a43d-49c2-9180-47ceba18599e/.system_generated/steps/14/content.md"
        if os.path.exists(local_path):
            with open(local_path, "r", encoding="utf-8") as f:
                html = f.read()
        else:
            raise e

    print("Parsing HTML...")
    parser = FathersParser()
    parser.feed(html)
    
    # Save the mapped tasks
    output_path = "/Users/raenard/Documents/devotion/research/task_map.json"
    os.makedirs(os.path.dirname(output_path), exist_ok=True)
    with open(output_path, "w", encoding="utf-8") as f:
        json.dump(parser.fathers, f, indent=2, ensure_ascii=False)
        
    print(f"Successfully generated task list with {len(parser.fathers)} fathers at {output_path}")

if __name__ == "__main__":
    main()
