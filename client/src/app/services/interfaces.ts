export interface ApiResult {
    page: number;
    results: any[];
    total_pages: number;
    total_results: number;
}

export interface EntryResult {
    id: number;
    entry_title: string;
    entry_overview: string;
    like_count: number;
}
