import os
import re
import json
import time
import urllib.request
from html.parser import HTMLParser

class ContentParser(HTMLParser):
    def __init__(self):
        super().__init__()
        self.sections = []
        self.current_section = None
        self.in_heading = False
        self.heading_tag = None
        self.heading_parts = []
        self.in_paragraph = False
        self.in_blockquote = False
        self.paragraph_parts = []
        self.in_anchor = False
        self.anchor_text_parts = []
        self.skip_all = False

    def handle_starttag(self, tag, attrs):
        if self.skip_all:
            return
        attrs_dict = dict(attrs)
        
        # Stop parsing when we hit the footer / about this page section
        if tag == 'div' and (attrs_dict.get('class') == 'pub' or attrs_dict.get('id') == 'ogdenville'):
            self.skip_all = True
            return
            
        if tag in ['h2', 'h3', 'h4']:
            self.in_heading = True
            self.heading_tag = tag
            self.heading_parts = []
        elif tag == 'p':
            # Check if this paragraph is advertising
            if attrs_dict.get('class') == 'h1a':
                # Sometimes a sub-summary, let's keep it
                pass
            self.in_paragraph = True
            self.paragraph_parts = []
        elif tag == 'blockquote':
            self.in_blockquote = True
            self.paragraph_parts = []
        elif tag == 'a' and (self.in_paragraph or self.in_blockquote or self.in_heading):
            self.in_anchor = True
            self.anchor_text_parts = []

    def handle_endtag(self, tag):
        if self.skip_all:
            return
        if tag in ['h2', 'h3', 'h4'] and self.in_heading:
            self.in_heading = False
            heading_text = "".join(self.heading_parts).strip()
            
            # Stop if we hit the "About this page" section
            if "about this page" in heading_text.lower():
                self.skip_all = True
                return
                
            if heading_text:
                self.current_section = {
                    "heading": heading_text,
                    "paragraphs": []
                }
                self.sections.append(self.current_section)
        elif tag == 'p' and self.in_paragraph:
            self.in_paragraph = False
            p_text = "".join(self.paragraph_parts).strip()
            # Clean up advertising/support message
            if p_text and "please help support the mission of new advent" not in p_text.lower():
                # If no heading was parsed yet, create a default one
                if not self.current_section:
                    self.current_section = {
                        "heading": "Preface",
                        "paragraphs": []
                    }
                    self.sections.append(self.current_section)
                self.current_section["paragraphs"].append(p_text)
        elif tag == 'blockquote' and self.in_blockquote:
            self.in_blockquote = False
            p_text = "".join(self.paragraph_parts).strip()
            if p_text:
                if not self.current_section:
                    self.current_section = {
                        "heading": "Preface",
                        "paragraphs": []
                    }
                    self.sections.append(self.current_section)
                self.current_section["paragraphs"].append(f"> {p_text}")
        elif tag == 'a' and self.in_anchor:
            self.in_anchor = False
            anchor_text = "".join(self.anchor_text_parts)
            if self.in_heading:
                self.heading_parts.append(anchor_text)
            else:
                self.paragraph_parts.append(anchor_text)

    def handle_data(self, data):
        if self.skip_all:
            return
        if self.in_anchor:
            self.anchor_text_parts.append(data)
        elif self.in_heading:
            self.heading_parts.append(data)
        elif self.in_paragraph or self.in_blockquote:
            self.paragraph_parts.append(data)

def clean_text(text):
    # Normalize whitespaces
    text = re.sub(r'\s+', ' ', text)
    # Remove HTML entities if any remaining
    text = text.replace('&nbsp;', ' ').replace('&#151;', '—').replace('&mdash;', '—').replace('&ndash;', '–')
    text = text.replace('&ldquo;', '"').replace('&rdquo;', '"').replace('&lsquo;', "'").replace('&rsquo;', "'")
    text = text.replace('&#145;', "'").replace('&#146;', "'").replace('&#147;', '"').replace('&#148;', '"')
    return text.strip()

def fetch_html(url):
    print(f"Fetching {url}...")
    time.sleep(3.0)  # Rate limit: 3 seconds delay for every download to not overload the site
    req = urllib.request.Request(url, headers={'User-Agent': 'Mozilla/5.0'})
    with urllib.request.urlopen(req) as response:
        return response.read().decode('utf-8', errors='ignore')

def parse_work_page(url):
    html = fetch_html(url)
    
    # Check if this is an index page (e.g. contains links to sub-books)
    # Search for links pointing to other sub-pages in /fathers/
    # For example, in 1201.htm, links are ../fathers/120101.htm, etc.
    # Exclude links back to index.html or other fathers
    # We can parse the links using a simple regex
    links = re.findall(r'href="(?:\.\./)?fathers/(\d+\.htm)"', html)
    if not links:
        # Check alternative link formats
        links = re.findall(r'href="(\d+\.htm)"', html)
        
    if links:
        # Multi-page work! Filter and reconstruct full urls
        # Deduplicate while preserving order
        unique_links = []
        for l in links:
            if l not in unique_links:
                unique_links.append(l)
        
        # Resolve subpages
        base_url = url.rsplit('/', 1)[0]
        subpage_urls = [f"{base_url}/{l}" for l in unique_links]
        print(f"Detected multi-page work. Found {len(subpage_urls)} subpages.")
        
        combined_sections = []
        for s_url in subpage_urls:
            try:
                sub_html = fetch_html(s_url)
                parser = ContentParser()
                parser.feed(sub_html)
                combined_sections.extend(parser.sections)
            except Exception as e:
                print(f"Error parsing subpage {s_url}: {e}")
        return combined_sections
    else:
        # Single page work
        parser = ContentParser()
        parser.feed(html)
        return parser.sections

def roman_numeral(n):
    val = [1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1]
    syb = ["M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"]
    roman_num = ''
    i = 0
    while n > 0:
        for _ in range(n // val[i]):
            roman_num += syb[i]
            n -= val[i]
        i += 1
    return roman_num

def generate_markdown(father_name, works):
    md = []
    md.append(f"# {father_name}\n")
    md.append(f"## {father_name}\n")
    md.append("Processed Church Father writings from New Advent.\n")
    
    for work in works:
        md.append(f"## {work['title']}")
        md.append("")
        
        sections = work.get("content_sections", [])
        for idx, sec in enumerate(sections):
            heading = clean_text(sec["heading"])
            # Remove chapter numbers if any to keep clean, or format as Chapter I: Heading
            # Standardize heading format to: #### Chapter I: Heading
            # First, check if the heading already contains chapter info
            heading_clean = re.sub(r'^(?:Chapter|Book|Section)\s+\d+[\.:\s]*', '', heading, flags=re.IGNORECASE).strip()
            heading_clean = re.sub(r'^\d+[\.:\s]*', '', heading_clean).strip()
            
            num = roman_numeral(idx + 1)
            md.append(f"#### Chapter {num}: {heading_clean}")
            md.append("")
            
            for p in sec["paragraphs"]:
                p_clean = clean_text(p)
                if p_clean:
                    md.append(p_clean)
                    md.append("")
        md.append("")
        
    return "\n".join(md)

def main():
    task_map_path = "/Users/raenard/Documents/devotion/research/task_map.json"
    resources_dir = "/Users/raenard/Documents/devotion/src/components/resources"
    
    if not os.path.exists(task_map_path):
        print("Task map not found. Please run task_builder.py first.")
        return

    while True:
        # Always reload the file to get the latest status
        with open(task_map_path, "r", encoding="utf-8") as f:
            fathers = json.load(f)
            
        # Find the next pending father
        target_father = None
        completed_count = 0
        total_count = len(fathers)
        
        for f in fathers:
            if f.get("status") == "completed":
                completed_count += 1
            elif f.get("status") == "pending" and target_father is None:
                target_father = f
                
        if not target_father:
            print(f"All fathers have been processed! ({completed_count}/{total_count} completed)")
            break
            
        print(f"\nProcessing Church Father: {target_father['name']} ({completed_count + 1}/{total_count})")
        
        processed_works = []
        for work in target_father["works"]:
            print(f"Downloading work: {work['title']} ({work['url']})")
            try:
                sections = parse_work_page(work["url"])
                processed_works.append({
                    "title": work["title"],
                    "content_sections": sections
                })
                work["status"] = "completed"
            except Exception as e:
                print(f"Failed to process work {work['title']}: {e}")
                work["status"] = "failed"
                
        # Generate MD filename
        safe_name = re.sub(r'[^a-z0-9]+', '_', target_father["name"].lower()).strip('_')
        md_filename = f"{safe_name}.md"
        md_filepath = os.path.join(resources_dir, md_filename)
        
        print(f"Formatting and saving to {md_filepath}...")
        md_content = generate_markdown(target_father["name"], processed_works)
        
        with open(md_filepath, "w", encoding="utf-8") as f_out:
            f_out.write(md_content)
            
        # Mark father as completed
        target_father["status"] = "completed"
        
        # Save the updated task map
        with open(task_map_path, "w", encoding="utf-8") as f_save:
            json.dump(fathers, f_save, indent=2, ensure_ascii=False)
            
        print(f"Finished processing {target_father['name']}. Markdown saved to {md_filepath}")

if __name__ == "__main__":
    main()
